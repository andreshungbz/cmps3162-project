// renderLoading displays loading text.
export function renderLoading(text = 'Loading...') {
  return `
    <div class="ui-loading">
      <div class="spinner"></div>
      <span>${text}</span>
    </div>
  `;
}

// renderError displays an error message from the server.
export function renderError(message = 'Something went wrong!') {
  if (!message) return '';

  return `
    <div class="ui-error">
      <strong>⚠ Error</strong>
      <p>${message}</p>
    </div>
  `;
}
