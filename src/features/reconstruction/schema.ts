import { z } from 'zod'

const vec3Schema = z.tuple([z.number(), z.number(), z.number()])

export const artifactSchema = z.object({
  caseId: z.string(),
  caseName: z.string(),
  version: z.string(),
  commit: z.string(),
  createdAt: z.string(),
  coordinateFrame: z.string(),
  uploads: z.array(
    z.object({
      fileName: z.string(),
      sizeBytes: z.number(),
      durationSeconds: z.number(),
      width: z.number(),
      height: z.number(),
      frameRate: z.number(),
      sha256: z.string(),
    }),
  ),
  points: z.array(
    z.object({
      id: z.string(),
      x: z.number(),
      y: z.number(),
      z: z.number(),
      color: z.array(z.number()).length(3),
      tags: z.array(z.string()).optional(),
    }),
  ),
  cameras: z.array(
    z.object({
      id: z.string(),
      sourceVideo: z.string(),
      timeSeconds: z.number(),
      position: vec3Schema.or(z.array(z.number()).length(3)),
      rotationEuler: vec3Schema.or(z.array(z.number()).length(3)),
      focalLengthPx: z.number(),
    }),
  ),
  vehicleTrack: z.array(
    z.object({
      timeSeconds: z.number(),
      x: z.number(),
      y: z.number(),
      z: z.number(),
      speedMps: z.number(),
    }),
  ),
  speed: z.object({
    method: z.string(),
    meanMps: z.number(),
    meanKph: z.number(),
    meanMph: z.number(),
    lowerKph: z.number(),
    upperKph: z.number(),
    confidence: z.number(),
    notes: z.array(z.string()),
  }),
  quality: z.object({
    mode: z.string(),
    inputVideos: z.number(),
    coordinateFrame: z.string(),
    toolchain: z.array(
      z.object({
        name: z.string(),
        path: z.string().optional(),
        version: z.string().optional(),
        status: z.string(),
      }),
    ),
    warnings: z.array(z.string()),
    reprojectionRmse: z.number().optional(),
  }),
  reportMarkdown: z.string(),
})
