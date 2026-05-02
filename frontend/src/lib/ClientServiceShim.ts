export async function ReorderClient(clientID: number, prevID: number, nextID: number) {
  // Call the Wails-generated binding if available; fallback will throw at runtime.
  // Use any to avoid TypeScript errors from generated bindings.
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  const svc = (window as any)?.go?.service?.ClientService;
  if (!svc || typeof svc.ReorderClient !== "function") {
    throw new Error("ReorderClient binding not available; regenerate wails bindings or run backend build");
  }
  return svc.ReorderClient(clientID, prevID, nextID);
}