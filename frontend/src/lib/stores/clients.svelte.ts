import {
  CreateClient,
  ListClients,
  ListClientsByZone,
  UpdateClient,
  UnregisterClient,
  ListInactiveClients,
  RestoreClient,
} from "wailsjs/go/service/ClientService";

import { ReorderClient as ReorderClientBinding } from "$lib/ClientServiceShim";
import type { model } from "wailsjs/go/models";

let clients = $state<model.Client[]>([]);
let loading = $state(false);

export function getClients() {
  return clients;
}

export function getClientsLoading() {
  return loading;
}

export async function fetchClients() {
  loading = true;
  try {
    clients = await ListClients();
  } finally {
    loading = false;
  }
}

export async function fetchClientsByZone(zoneId: number) {
  loading = true;
  try {
    clients = (await ListClientsByZone(zoneId)) ?? [];
  } finally {
    loading = false;
  }
}

export async function createClient(
  name: string,
  address: string,
  zoneId: number,
) {
  loading = true;
  try {
    await CreateClient(name, address, zoneId);
    await fetchClients();
  } finally {
    loading = false;
  }
}

export async function updateClient(
  id: number,
  name: string,
  address: string,
  zoneId: number,
) {
  loading = true;
  try {
    await UpdateClient(id, name, address, zoneId);
    await fetchClients();
  } finally {
    loading = false;
  }
}

export async function unregisterClient(id: number) {
  loading = true;
  try {
    await UnregisterClient(id);
    await fetchClients();
  } finally {
    loading = false;
  }
}

export async function listInactiveClients() {
  try {
    return await ListInactiveClients();
  } catch (err) {
    console.error("Error listing inactive clients:", err);
    return [];
  }
}

export async function restoreClient(id: number) {
  loading = true;
  try {
    await RestoreClient(id);
    await fetchClients();
  } finally {
    loading = false;
  }
}

// Reorder a client within its zone. prevID or nextID may be 0 to indicate start/end.
export async function reorderClient(clientID: number, prevID: number, nextID: number) {
  loading = true;
  try {
    await ReorderClientBinding(clientID, prevID, nextID);
    // Refresh clients. Caller should decide whether to call fetchClientsByZone for a focused zone.
    await fetchClients();
  } finally {
    loading = false;
  }
}
