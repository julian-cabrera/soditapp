import {
  CreateClient,
  ListClients,
  ListClientsByZone,
  UpdateClient,
  UnregisterClient,
} from "wailsjs/go/service/ClientService";
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
    clients = await ListClientsByZone(zoneId);
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
