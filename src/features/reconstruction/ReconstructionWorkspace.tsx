import { lazy, Suspense, useEffect, useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import {
  BadgeDollarSign,
  Download,
  FileArchive,
  Map,
  Play,
  RefreshCw,
  Scale,
  Star,
  Upload,
  Video,
} from "lucide-react";
import {
  createCase,
  getApiBase,
  getArtifact,
  getCase,
  listTools,
  setApiBase,
  type Artifact,
  type CaseSummary,
} from "../../api/client";
import { APP_BUILT_AT, APP_COMMIT, APP_VERSION } from "../../generated/version";
import { loadLastArtifact, saveLastArtifact } from "./storage";
import { sampleArtifact } from "./sampleArtifact";

const SceneViewer = lazy(() => import("./SceneViewer"));

const repoUrl = "https://github.com/baditaflorin/accident-reconstructor";
const paypalUrl = "https://www.paypal.com/paypalme/florinbadita";

export function ReconstructionWorkspace() {
  const [apiBaseInput, setApiBaseInput] = useState(getApiBase);
  const [apiBase, setCurrentApiBase] = useState(getApiBase);
  const [caseName, setCaseName] = useState("Intersection reconstruction");
  const [scaleMeters, setScaleMeters] = useState(10);
  const [files, setFiles] = useState<File[]>([]);
  const [dragActive, setDragActive] = useState(false);
  const [summary, setSummary] = useState<CaseSummary | null>(null);
  const [artifact, setArtifact] = useState<Artifact | null>(null);
  const [toast, setToast] = useState("");

  useEffect(() => {
    loadLastArtifact().then((cached) => {
      if (cached) {
        setArtifact(cached);
      }
    });
  }, []);

  const toolsQuery = useQuery({
    queryKey: ["tools", apiBase],
    queryFn: () => listTools(apiBase),
  });

  const uploadMutation = useMutation({
    mutationFn: async () => {
      if (files.length === 0) {
        throw new Error("Add at least one dashcam or phone video first.");
      }
      const created = await createCase({
        apiBase,
        caseName,
        scaleMeters,
        files,
      });
      setSummary(created);
      const completed = await pollUntilComplete(
        apiBase,
        created.id,
        setSummary,
      );
      const nextArtifact = await getArtifact(apiBase, completed.id);
      setArtifact(nextArtifact);
      await saveLastArtifact(nextArtifact);
      return nextArtifact;
    },
    onError(error) {
      setToast(
        error instanceof Error ? error.message : "Reconstruction failed.",
      );
    },
  });

  const metrics = useMemo(() => {
    if (!artifact) {
      return [
        ["Speed", "Waiting"],
        ["Confidence", "Waiting"],
        ["Sparse Points", "0"],
        ["Videos", String(files.length)],
      ];
    }
    return [
      ["Speed", `${artifact.speed.meanKph.toFixed(1)} km/h`],
      ["Confidence", `${Math.round(artifact.speed.confidence * 100)}%`],
      ["Sparse Points", String(artifact.points.length)],
      ["Videos", String(artifact.uploads.length)],
    ];
  }, [artifact, files.length]);

  function applyApiBase() {
    setApiBase(apiBaseInput);
    setCurrentApiBase(apiBaseInput.replace(/\/$/, ""));
  }

  function handleFiles(nextFiles: FileList | File[]) {
    setFiles(Array.from(nextFiles));
    setToast("");
  }

  function downloadArtifact() {
    if (!artifact) {
      return;
    }
    const blob = new Blob([JSON.stringify(artifact, null, 2)], {
      type: "application/json",
    });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${artifact.caseId}-reconstruction.json`;
    link.click();
    URL.revokeObjectURL(url);
  }

  return (
    <div className="app-shell">
      <header className="topbar">
        <div className="brand">
          <div className="brand-mark" aria-hidden="true">
            <Scale size={23} />
          </div>
          <div>
            <h1>Accident Reconstructor</h1>
            <p>
              3D crash scene reconstruction, speed ranges, and evidence bundles.
            </p>
          </div>
        </div>
        <nav className="top-actions" aria-label="Project links">
          <a
            className="btn btn-dark"
            href={repoUrl}
            target="_blank"
            rel="noreferrer"
            title="Star the GitHub repository"
          >
            <Star size={17} />
            Star on GitHub
          </a>
          <a
            className="btn"
            href={paypalUrl}
            target="_blank"
            rel="noreferrer"
            title="Support the project on PayPal"
          >
            <BadgeDollarSign size={17} />
            Support
          </a>
        </nav>
      </header>

      <main className="workspace">
        <aside className="side-panel" aria-label="Reconstruction controls">
          <section className="section">
            <h2>Case Intake</h2>
            <div className="field">
              <label htmlFor="case-name">Case name</label>
              <input
                id="case-name"
                className="input"
                value={caseName}
                onChange={(event) => setCaseName(event.target.value)}
              />
            </div>
            <div className="field">
              <label htmlFor="scale-meters">
                Measured road reference in meters
              </label>
              <input
                id="scale-meters"
                className="input"
                min="1"
                step="0.5"
                type="number"
                value={scaleMeters}
                onChange={(event) => setScaleMeters(Number(event.target.value))}
              />
            </div>
            <label
              className="drop-zone"
              data-active={dragActive}
              onDragEnter={(event) => {
                event.preventDefault();
                setDragActive(true);
              }}
              onDragOver={(event) => event.preventDefault()}
              onDragLeave={() => setDragActive(false)}
              onDrop={(event) => {
                event.preventDefault();
                setDragActive(false);
                handleFiles(event.dataTransfer.files);
              }}
            >
              <input
                className="sr-only"
                type="file"
                multiple
                accept="video/*"
                onChange={(event) => handleFiles(event.target.files ?? [])}
              />
              <span>
                <Upload size={26} aria-hidden="true" />
                <strong>Drop dashcam and phone videos</strong>
                <span className="hint">
                  or click to choose files for backend reconstruction
                </span>
              </span>
            </label>
            <div className="file-list" aria-live="polite">
              {files.map((file) => (
                <div className="file-row" key={`${file.name}-${file.size}`}>
                  <span>
                    <Video size={15} aria-hidden="true" /> {file.name}
                  </span>
                  <span className="badge">{formatBytes(file.size)}</span>
                </div>
              ))}
            </div>
          </section>

          <section className="section">
            <h2>Runtime API</h2>
            <div className="field">
              <label htmlFor="api-base">Backend URL</label>
              <input
                id="api-base"
                className="input"
                value={apiBaseInput}
                onChange={(event) => setApiBaseInput(event.target.value)}
              />
            </div>
            <div className="button-row">
              <button className="btn" type="button" onClick={applyApiBase}>
                <RefreshCw size={16} />
                Use API
              </button>
              <button
                className="btn"
                type="button"
                onClick={() => setArtifact(sampleArtifact)}
              >
                <Map size={16} />
                Load Sample
              </button>
            </div>
            <div className="file-list">
              {(toolsQuery.data ?? []).slice(0, 6).map((tool) => (
                <div className="tool-row" key={tool.name}>
                  <span>{tool.name}</span>
                  <span
                    className="badge"
                    data-tone={tool.status === "available" ? undefined : "warn"}
                  >
                    {tool.status}
                  </span>
                </div>
              ))}
              {toolsQuery.isError && (
                <div className="warning-row">
                  <span>
                    Backend unavailable. The sample viewer still works on Pages.
                  </span>
                </div>
              )}
            </div>
          </section>

          <section className="section">
            <h2>Reconstruction</h2>
            {summary && (
              <>
                <div
                  className="progress"
                  aria-label={`Progress ${summary.progress}%`}
                >
                  <span style={{ width: `${summary.progress}%` }} />
                </div>
                <p className="hint">{summary.message}</p>
              </>
            )}
            <div className="button-row">
              <button
                className="btn btn-primary"
                type="button"
                disabled={uploadMutation.isPending || files.length === 0}
                onClick={() => uploadMutation.mutate()}
              >
                <Play size={16} />
                Run
              </button>
              <button
                className="btn"
                type="button"
                disabled={!artifact}
                onClick={downloadArtifact}
              >
                <Download size={16} />
                JSON
              </button>
              {summary?.status === "complete" && summary.artifactUrl && (
                <a
                  className="btn"
                  href={`${apiBase}${summary.artifactUrl}`}
                  target="_blank"
                  rel="noreferrer"
                >
                  <FileArchive size={16} />
                  API Artifact
                </a>
              )}
            </div>
            {toast && <div className="toast">{toast}</div>}
          </section>

          {artifact && (
            <section className="section">
              <h2>Evidence Report</h2>
              <div className="report">{artifact.reportMarkdown}</div>
            </section>
          )}
        </aside>

        <section
          className="main-panel"
          aria-label="Reconstruction visualization"
        >
          <Suspense
            fallback={
              <div className="viewer-shell">
                <div className="canvas-empty">Loading 3D scene...</div>
              </div>
            }
          >
            <SceneViewer artifact={artifact} />
          </Suspense>
          <div className="evidence-strip">
            {metrics.map(([label, value]) => (
              <div className="metric" key={label}>
                <span>{label}</span>
                <strong>{value}</strong>
              </div>
            ))}
          </div>
        </section>
      </main>

      <footer className="footer">
        <span>
          Version {APP_VERSION} · commit {APP_COMMIT} · built{" "}
          {formatBuiltAt(APP_BUILT_AT)}
        </span>
        <span className="footer-links">
          <a href={repoUrl}>
            https://github.com/baditaflorin/accident-reconstructor
          </a>
          <a href={paypalUrl}>https://www.paypal.com/paypalme/florinbadita</a>
        </span>
      </footer>
    </div>
  );
}

async function pollUntilComplete(
  apiBase: string,
  caseId: string,
  onUpdate: (summary: CaseSummary) => void,
) {
  for (let attempt = 0; attempt < 120; attempt += 1) {
    const next = await getCase(apiBase, caseId);
    onUpdate(next);
    if (next.status === "complete") {
      return next;
    }
    if (next.status === "failed") {
      throw new Error(next.error?.message ?? "Reconstruction failed.");
    }
    await new Promise((resolve) => setTimeout(resolve, 1_000));
  }
  throw new Error("Timed out waiting for reconstruction.");
}

function formatBytes(value: number) {
  if (value < 1024 * 1024) {
    return `${(value / 1024).toFixed(1)} KB`;
  }
  return `${(value / (1024 * 1024)).toFixed(1)} MB`;
}

function formatBuiltAt(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toISOString().slice(0, 10);
}
