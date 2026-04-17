/**
 * Future helper:
 * Convert backend job → UI model
 */
export function normalizeJob(job) {
  return job;
}

/**
 * Future helper:
 * Map status → UI badge color
 */
export function getStatusClass(status) {
  switch (status) {
    case 'done':
      return 'success';
    case 'processing':
      return 'warning';
    case 'failed':
      return 'error';
    default:
      return 'pending';
  }
}
