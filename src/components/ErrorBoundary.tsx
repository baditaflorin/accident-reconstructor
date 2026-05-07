import { Component, type ErrorInfo, type ReactNode } from 'react'

type Props = {
  children: ReactNode
}

type State = {
  hasError: boolean
}

export class ErrorBoundary extends Component<Props, State> {
  state: State = { hasError: false }

  static getDerivedStateFromError(): State {
    return { hasError: true }
  }

  componentDidCatch(error: Error, info: ErrorInfo) {
    if (import.meta.env.DEV) {
      console.error(error, info)
    }
  }

  render() {
    if (this.state.hasError) {
      return (
        <main className="app-shell">
          <div className="canvas-empty">
            <div>
              <h1>Reconstruction workspace crashed</h1>
              <p>Reload the page and try again. Local drafts stay in browser storage.</p>
            </div>
          </div>
        </main>
      )
    }

    return this.props.children
  }
}
