import { openDB } from 'idb'
import type { Artifact } from '../../api/client'
import { artifactSchema } from './schema'

const DB_NAME = 'accident-reconstructor'
const STORE = 'artifacts'
const LAST_KEY = 'last-artifact'

async function db() {
  return openDB(DB_NAME, 1, {
    upgrade(database) {
      database.createObjectStore(STORE)
    },
  })
}

export async function saveLastArtifact(artifact: Artifact) {
  const database = await db()
  await database.put(STORE, artifact, LAST_KEY)
}

export async function loadLastArtifact() {
  const database = await db()
  const value = await database.get(STORE, LAST_KEY)
  const parsed = artifactSchema.safeParse(value)
  return parsed.success ? (parsed.data as Artifact) : null
}
