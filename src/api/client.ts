import createClient from "openapi-fetch";
import type { components, paths } from "./schema";

export type Artifact = components["schemas"]["Artifact"];
export type CaseSummary = components["schemas"]["CaseSummary"];
export type ToolStatus = components["schemas"]["ToolStatus"];

const API_STORAGE_KEY = "accident-reconstructor:api-base";
const DEFAULT_API_BASE =
  import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export class ApiError extends Error {
  code: string;

  constructor(message: string, code = "api_error") {
    super(message);
    this.name = "ApiError";
    this.code = code;
  }
}

export function getApiBase() {
  return window.localStorage.getItem(API_STORAGE_KEY) || DEFAULT_API_BASE;
}

export function setApiBase(value: string) {
  window.localStorage.setItem(API_STORAGE_KEY, value.replace(/\/$/, ""));
}

export function apiClient(baseUrl = getApiBase()) {
  return createClient<paths>({ baseUrl });
}

function toApiError(error: unknown, fallback: string) {
  if (error && typeof error === "object") {
    const maybe = error as { message?: unknown; code?: unknown };
    return new ApiError(
      typeof maybe.message === "string" ? maybe.message : fallback,
      typeof maybe.code === "string" ? maybe.code : "api_error",
    );
  }
  return new ApiError(fallback);
}

export async function listTools(apiBase: string) {
  const { data, error } = await apiClient(apiBase).GET("/api/v1/tools");
  if (error) {
    throw toApiError(error, "Unable to read backend toolchain status.");
  }
  return data ?? [];
}

export async function createCase(input: {
  apiBase: string;
  caseName: string;
  scaleMeters: number;
  files: File[];
}) {
  const body = new FormData();
  body.set("case_name", input.caseName);
  body.set("scale_meters", String(input.scaleMeters));
  for (const file of input.files) {
    body.append("videos", file);
  }

  const { data, error } = await apiClient(input.apiBase).POST("/api/v1/cases", {
    body: body as never,
    bodySerializer: (payload) => payload as unknown as BodyInit,
  });
  if (error) {
    throw toApiError(error, "Unable to create reconstruction case.");
  }
  if (!data) {
    throw new ApiError("Backend returned an empty create-case response.");
  }
  return data.case;
}

export async function getCase(apiBase: string, caseId: string) {
  const { data, error } = await apiClient(apiBase).GET(
    "/api/v1/cases/{caseId}",
    {
      params: { path: { caseId } },
    },
  );
  if (error) {
    throw toApiError(error, "Unable to read reconstruction case.");
  }
  if (!data) {
    throw new ApiError("Backend returned an empty case response.");
  }
  return data;
}

export async function getArtifact(apiBase: string, caseId: string) {
  const { data, error } = await apiClient(apiBase).GET(
    "/api/v1/cases/{caseId}/artifact",
    {
      params: { path: { caseId } },
    },
  );
  if (error) {
    throw toApiError(error, "Unable to read reconstruction artifact.");
  }
  if (!data) {
    throw new ApiError("Backend returned an empty artifact response.");
  }
  return data;
}
