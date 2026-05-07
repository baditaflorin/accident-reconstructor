import type { Artifact } from '../../api/client'

export const sampleArtifact: Artifact = {
  caseId: 'sample-case',
  caseName: 'Sample two-camera intersection case',
  version: '0.1.0',
  commit: 'sample',
  createdAt: new Date('2026-05-08T00:00:00Z').toISOString(),
  coordinateFrame: 'local_enu_meters',
  uploads: [
    {
      fileName: 'dashcam-eastbound.mp4',
      sizeBytes: 18_400_000,
      durationSeconds: 6.2,
      width: 1920,
      height: 1080,
      frameRate: 29.97,
      sha256: 'sample-dashcam',
    },
    {
      fileName: 'bystander-north-corner.mov',
      sizeBytes: 24_900_000,
      durationSeconds: 6.6,
      width: 1280,
      height: 720,
      frameRate: 30,
      sha256: 'sample-bystander',
    },
  ],
  points: Array.from({ length: 42 }, (_, index) => {
    const lane = (index % 3) - 1
    const step = Math.floor(index / 3)
    return {
      id: `sample-point-${index}`,
      x: -24 + step * 3.6,
      y: 0,
      z: lane * 3.4,
      color: lane === 0 ? [242, 205, 71] : [176, 184, 188],
      tags: ['road'],
    }
  }),
  cameras: Array.from({ length: 10 }, (_, index) => ({
    id: `sample-camera-${index}`,
    sourceVideo: index < 5 ? 'dashcam-eastbound.mp4' : 'bystander-north-corner.mov',
    timeSeconds: (index % 5) * 1.2,
    position: [-24 + (index % 5) * 10, 2.4, index < 5 ? -8 : 9],
    rotationEuler: [-8, 0, 0],
    focalLengthPx: 900,
  })),
  vehicleTrack: Array.from({ length: 18 }, (_, index) => {
    const progress = index / 17
    return {
      timeSeconds: progress * 6.2,
      x: -24 + progress * 48,
      y: 0,
      z: Math.sin(progress * Math.PI * 2) * 0.8,
      speedMps: 7.74,
    }
  }),
  speed: {
    method: 'track-distance-over-time',
    meanMps: 7.74,
    meanKph: 27.86,
    meanMph: 17.31,
    lowerKph: 21.91,
    upperKph: 33.8,
    confidence: 0.64,
    notes: [
      'Multiple videos improve camera path constraints.',
      'Use a measured road reference to improve scale accuracy.',
    ],
  },
  quality: {
    mode: 'sample-artifact',
    inputVideos: 2,
    coordinateFrame: 'local_enu_meters',
    toolchain: [
      { name: 'colmap', status: 'available' },
      { name: 'ffmpeg', status: 'available' },
      { name: 'gdalinfo', status: 'available' },
      { name: 'ollama', status: 'optional' },
    ],
    warnings: ['Sample data for UI demonstration only.'],
  },
  reportMarkdown:
    '# Accident Reconstruction Report\n\nSample two-camera case with sparse road points, camera poses, speed range, and reconstruction warnings.',
}
