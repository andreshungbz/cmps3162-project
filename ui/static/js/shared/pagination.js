// getPaginationItems constructs an array of proper pagination values
// based on the current page and total pages.
export function getPaginationItems(current, total) {
  const pages = [];

  // 1 is always present
  pages.push(1);

  if (current > 3) pages.push('...');

  // determine middle pagination buttons shown
  const start = Math.max(2, current - 2);
  const end = Math.min(total - 1, current + 2);
  for (let i = start; i <= end; i++) {
    pages.push(i);
  }

  if (current < total - 2) pages.push('...');

  // last page is always present (if more than 1 page)
  if (total > 1) pages.push(total);

  return pages;
}
