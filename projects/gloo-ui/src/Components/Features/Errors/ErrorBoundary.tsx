import * as React from 'react';
import { Link } from 'react-router-dom';
interface Props {
  children?: React.ReactNode;
  noRedirect?: boolean;
}

export class ErrorBoundary extends React.Component<Props, any> {
  state = {
    hasError: false
  };
  static getDerivedStateFromError(error: Error) {
    return { hasError: true };
  }
  componentDidCatch(error: Error, info: object) {
    console.error('Error Boundary caught an error', error, info);
  }

  render() {
    if (this.state.hasError) {
      return (
        <div>
          <div>Something went wrong.</div>

          {!this.props.noRedirect && (
            <div>
              <Link to='/catalog' style={{ textDecoration: 'none' }}>
                Click here
              </Link>
              to go back to the Catalog page
            </div>
          )}
        </div>
      );
    }

    return this.props.children;
  }
}
