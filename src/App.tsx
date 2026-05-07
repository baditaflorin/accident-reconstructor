import { ErrorBoundary } from './components/ErrorBoundary'
import { ReconstructionWorkspace } from './features/reconstruction/ReconstructionWorkspace'

function App() {
  return (
    <ErrorBoundary>
      <ReconstructionWorkspace />
    </ErrorBoundary>
  )
}

export default App
