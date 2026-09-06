import {
  GetDeliverySheets,
  ConfirmDelivery,
} from "wailsjs/go/service/ClientProductService";
import { service, type dto } from "wailsjs/go/models";
import { formatISO, formatDate } from "$lib/utils/date.js";

export interface ClientGroup {
  clientId: number;
  clientName: string;
  clientAddress: string;
  clientBalance: number;
  entries: dto.DeliveryEntryDTO[];
}

export function groupEntriesByClient(
  entries: dto.DeliveryEntryDTO[],
): ClientGroup[] {
  const groups: ClientGroup[] = [];
  let lastClientId = -1;
  for (const entry of entries) {
    if (entry.client_id === lastClientId) {
      groups.at(-1)!.entries.push(entry);
    } else {
      groups.push({
        clientId: entry.client_id,
        clientName: entry.client_name,
        clientAddress: entry.client_address,
        clientBalance: entry.client_balance,
        entries: [entry],
      });
      lastClientId = entry.client_id;
    }
  }
  return groups;
}

let sheets = $state<dto.DeliverySheetDTO[]>([]);
let loading = $state(false);
let currentPage = $state(0);
let entryEdits = $state<Map<string, { charged: number; credit: number }>>(
  new Map(),
);
let selectedZoneId = $state<number | undefined>(undefined);
let selectedDate = $state<Date>(new Date());

function getWeekStart(date: Date): Date {
  const d = new Date(date);
  const day = d.getDay();
  const diff = (day === 0 ? -6 : 1) - day;
  d.setDate(d.getDate() + diff);
  d.setHours(0, 0, 0, 0);
  return d;
}

function getWeekEnd(date: Date): Date {
  const start = getWeekStart(date);
  const end = new Date(start);
  end.setDate(end.getDate() + 6);
  return end;
}

export function getSheets() {
  return sheets;
}

export function getDeliveryLoading() {
  return loading;
}

export function getCurrentPage() {
  return currentPage;
}

export function getEntryEdits() {
  return entryEdits;
}

export function getSelectedZoneId() {
  return selectedZoneId;
}

export function getSelectedDate() {
  return selectedDate;
}

export function getWeekStartDate(): Date {
  return getWeekStart(selectedDate);
}

export function getWeekEndDate(): Date {
  return getWeekEnd(selectedDate);
}

export function setSelectedDate(date: Date) {
  selectedDate = date;
}

export async function fetchDeliverySheets(zoneId: number) {
  selectedZoneId = zoneId;
  loading = true;
  try {
    sheets = await GetDeliverySheets(zoneId);
    currentPage = 0;
    const edits = new Map<string, { charged: number; credit: number }>();
    for (const sheet of sheets) {
      for (const entry of sheet.entries) {
        const key = `${entry.client_id}-${entry.product_id}`;
        if (!edits.has(key)) {
          edits.set(key, { charged: 0, credit: 0 });
        }
      }
    }
    entryEdits = edits;
  } catch {
    sheets = [];
    entryEdits = new Map();
  } finally {
    loading = false;
  }
}

export function setCurrentPage(page: number) {
  currentPage = page;
}

export function updateEntryEdit(
  clientId: number,
  productId: number,
  field: "charged" | "credit",
  value: number,
) {
  const key = `${clientId}-${productId}`;
  const current = entryEdits.get(key) ?? { charged: 0, credit: 0 };
  current[field] = Math.max(0, value);
  entryEdits.set(key, current);
  entryEdits = new Map(entryEdits);
}

export function formatDateShort(date: Date): string {
  return formatDate(date);
}

export async function confirmDelivery() {
  const deliveryDate = formatISO(selectedDate);
  const adjustments: Array<{
    client_id: number;
    product_id: number;
    delta: number;
    charged_units: number;
    delivery_date: string;
  }> = [];

  for (const [key, edit] of entryEdits) {
    if (edit.charged > 0 || edit.credit > 0) {
      const [clientId, productId] = key.split("-").map(Number);
      adjustments.push({
        client_id: clientId,
        product_id: productId,
        delta: -edit.credit,
        charged_units: edit.charged,
        delivery_date: deliveryDate,
      });
    }
  }
  if (adjustments.length === 0) return;

  loading = true;
  try {
    await ConfirmDelivery(
      adjustments as unknown as service.BalanceAdjustment[],
    );
    if (selectedZoneId !== undefined) {
      await fetchDeliverySheets(selectedZoneId);
    }
  } finally {
    loading = false;
  }
}
