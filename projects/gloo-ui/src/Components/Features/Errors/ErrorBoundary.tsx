import * as React from 'react';

type State = {
  hasError: boolean;
  error: Error | null;
};

type ErrorBoundProps = {
  fallback: React.ReactNode;
};

// Error boundaries currently have to be classes.
export class ErrorBoundary extends React.Component<ErrorBoundProps, State> {
  state = { hasError: false, error: null };
  static getDerivedStateFromError(error: Error) {
    return {
      hasError: true,
      error
    };
  }
  componentDidCatch(error: Error, info: object) {
    console.error('Error Boundary caught an error', error, info);
  }

  render() {
    if (this.state.hasError) {
      return this.props.fallback;
    }
    return this.props.children;
  }
}
