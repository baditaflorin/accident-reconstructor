import { useEffect, useRef } from 'react'
import * as THREE from 'three'
import type { Artifact } from '../../api/client'

type Props = {
  artifact: Artifact | null
}

export default function SceneViewer({ artifact }: Props) {
  const hostRef = useRef<HTMLDivElement | null>(null)

  useEffect(() => {
    const host = hostRef.current
    if (!host || !artifact) {
      return
    }

    const scene = new THREE.Scene()
    scene.background = new THREE.Color('#0e171b')

    const camera = new THREE.PerspectiveCamera(52, 1, 0.1, 1000)
    camera.position.set(18, 22, 34)
    camera.lookAt(0, 0, 0)

    const renderer = new THREE.WebGLRenderer({ antialias: true })
    renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2))
    host.appendChild(renderer.domElement)

    const ambient = new THREE.AmbientLight('#ffffff', 0.7)
    scene.add(ambient)
    const sun = new THREE.DirectionalLight('#ffffff', 1.2)
    sun.position.set(18, 26, 14)
    scene.add(sun)

    const grid = new THREE.GridHelper(80, 40, '#59656b', '#243036')
    scene.add(grid)

    const road = new THREE.Mesh(
      new THREE.BoxGeometry(80, 0.06, 8.4),
      new THREE.MeshStandardMaterial({ color: '#2d3438', roughness: 0.92 }),
    )
    road.position.y = -0.04
    scene.add(road)

    for (const point of artifact.points) {
      const color = new THREE.Color(
        `rgb(${point.color[0] ?? 220}, ${point.color[1] ?? 220}, ${point.color[2] ?? 220})`,
      )
      const marker = new THREE.Mesh(
        new THREE.SphereGeometry(point.tags?.includes('vehicle') ? 0.28 : 0.14, 12, 12),
        new THREE.MeshStandardMaterial({ color }),
      )
      marker.position.set(point.x, point.y + 0.1, point.z)
      scene.add(marker)
    }

    const trackPoints = artifact.vehicleTrack.map((point) => new THREE.Vector3(point.x, point.y + 0.24, point.z))
    const track = new THREE.Line(
      new THREE.BufferGeometry().setFromPoints(trackPoints),
      new THREE.LineBasicMaterial({ color: '#5cc8ff' }),
    )
    scene.add(track)

    const vehicle = new THREE.Mesh(
      new THREE.BoxGeometry(4.2, 1.5, 2),
      new THREE.MeshStandardMaterial({ color: '#e03a2f', metalness: 0.1, roughness: 0.55 }),
    )
    scene.add(vehicle)

    for (const pose of artifact.cameras) {
      const cameraMarker = new THREE.Mesh(
        new THREE.ConeGeometry(0.45, 1, 4),
        new THREE.MeshStandardMaterial({ color: '#f5f0d4' }),
      )
      cameraMarker.position.set(pose.position[0], pose.position[1], pose.position[2])
      cameraMarker.rotation.x = Math.PI / 2
      scene.add(cameraMarker)
    }

    const resize = () => {
      const rect = host.getBoundingClientRect()
      const width = Math.max(320, rect.width)
      const height = Math.max(360, rect.height)
      renderer.setSize(width, height, false)
      camera.aspect = width / height
      camera.updateProjectionMatrix()
    }
    const observer = new ResizeObserver(resize)
    observer.observe(host)
    resize()

    let frame = 0
    let raf = 0
    const animate = () => {
      frame += 1
      const trackIndex = Math.floor((frame / 7) % artifact.vehicleTrack.length)
      const position = artifact.vehicleTrack[trackIndex]
      if (position) {
        vehicle.position.set(position.x, position.y + 0.82, position.z)
      }
      scene.rotation.y = Math.sin(frame / 360) * 0.08
      renderer.render(scene, camera)
      raf = requestAnimationFrame(animate)
    }
    animate()

    return () => {
      cancelAnimationFrame(raf)
      observer.disconnect()
      host.removeChild(renderer.domElement)
      renderer.dispose()
      scene.traverse((object) => {
        if (object instanceof THREE.Mesh) {
          object.geometry.dispose()
          const materials = Array.isArray(object.material) ? object.material : [object.material]
          for (const material of materials) {
            material.dispose()
          }
        }
      })
    }
  }, [artifact])

  if (!artifact) {
    return (
      <div className="viewer-shell">
        <div className="canvas-empty">
          <div>
            <h2>Drop videos or load the sample case</h2>
            <p>The reconstruction scene will appear here with camera paths, sparse points, and vehicle track playback.</p>
          </div>
        </div>
      </div>
    )
  }

  return <div ref={hostRef} className="viewer-shell" aria-label="3D reconstruction viewer" />
}
