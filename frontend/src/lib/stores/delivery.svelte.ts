import {
  GetDeliverySheets,
  ConfirmDelivery,
} from "wailsjs/go/service/ClientProductService";
import { service, type dto } from "wailsjs/go/models";

let sheets = $state<dto.DeliverySheetDTO[]>([]);
let loading = $state(false);
let currentPage = $state(0);
let entryEdits = $state<Map<string, { charged: number; credit: number }>>(
  new Map(),
);
let selectedZoneId = $state<number | undefined>(undefined);

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

export async function confirmDelivery() {
  const adjustments: service.BalanceAdjustment[] = [];
  for (const [key, edit] of entryEdits) {
    if (edit.credit > 0) {
      const [clientId, productId] = key.split("-").map(Number);
      adjustments.push(
        new service.BalanceAdjustment({
          client_id: clientId,
          product_id: productId,
          delta: -edit.credit,
        }),
      );
    }
  }
  if (adjustments.length === 0) return;

  loading = true;
  try {
    await ConfirmDelivery(adjustments);
    if (selectedZoneId !== undefined) {
      await fetchDeliverySheets(selectedZoneId);
    }
  } finally {
    loading = false;
  }
}
