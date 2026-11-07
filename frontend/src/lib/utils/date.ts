/**
 * Centralized date utilities for SoditApp.
 *
 * Display format: dd-MM-yyyy (user-facing)
 * ISO format:    yyyy-MM-dd (backend communication)
 */

/** Pad a number to 2 digits. */
function pad2(n: number): string {
  return String(n).padStart(2, "0");
}

/** Format a Date to dd-MM-yyyy for display. */
export function formatDate(date: Date): string {
  const d = pad2(date.getDate());
  const m = pad2(date.getMonth() + 1);
  const y = date.getFullYear();
  return `${d}-${m}-${y}`;
}

/** Format a Date to yyyy-MM-dd for the Go backend. */
export function formatISO(date: Date): string {
  const y = date.getFullYear();
  const m = pad2(date.getMonth() + 1);
  const d = pad2(date.getDate());
  return `${y}-${m}-${d}`;
}

/** Convert an ISO date string (yyyy-MM-dd) to display format (dd-MM-yyyy). */
export function isoToDisplay(isoDate: string): string {
  const [year, month, day] = isoDate.split("-");
  return `${day}-${month}-${year}`;
}

/** Convert a display string (dd-MM-yyyy) back to ISO (yyyy-MM-dd). */
export function displayToISO(displayDate: string): string {
  const [day, month, year] = displayDate.split("-");
  return `${year}-${month}-${day}`;
}

/** Default end date: today in dd-MM-yyyy. */
export function defaultEndDate(): string {
  return formatDate(new Date());
}

/** Default start date: N days ago in dd-MM-yyyy. */
export function defaultStartDate(daysBack = 30): string {
  const d = new Date();
  d.setDate(d.getDate() - daysBack);
  return formatDate(d);
}
