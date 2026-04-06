import {
  CreateZone,
  ListZones,
  UpdateZone,
  DeleteZone,
} from "wailsjs/go/service/ZoneService";
import type { model } from "wailsjs/go/models";

let zones = $state<model.Zone[]>([]);
let loading = $state(false);

export function getZones() {
  return zones;
}

export function getZonesLoading() {
  return loading;
}

export async function fetchZones() {
  loading = true;
  try {
    zones = await ListZones();
  } finally {
    loading = false;
  }
}

export async function createZone(name: string) {
  loading = true;
  try {
    await CreateZone(name);
    await fetchZones();
  } finally {
    loading = false;
  }
}

export async function updateZone(id: number, name: string) {
  loading = true;
  try {
    await UpdateZone(id, name);
    await fetchZones();
  } finally {
    loading = false;
  }
}

export async function deleteZone(id: number) {
  loading = true;
  try {
    await DeleteZone(id);
    await fetchZones();
  } finally {
    loading = false;
  }
}
