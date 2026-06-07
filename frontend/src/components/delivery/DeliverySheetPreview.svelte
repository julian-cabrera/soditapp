<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import { groupEntriesByClient } from "$lib/stores/delivery.svelte.js";
  import type { dto } from "wailsjs/go/models";

  let {
    sheets,
    onBack,
  }: {
    sheets: dto.DeliverySheetDTO[];
    onBack: () => void;
  } = $props();

  const dayColumns = [1, 2, 3, 4];

  function formatDate(value: string | Date): string {
    const date = typeof value === "string" ? new Date(value) : value;
    const d = String(date.getDate()).padStart(2, "0");
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const y = date.getFullYear();
    return `${d}/${m}/${y}`;
  }
</script>

<div class="print-hidden mb-4 flex items-center justify-between">
  <Button variant="outline" onclick={onBack}>Volver</Button>
  <Button onclick={() => window.print()}>Imprimir</Button>
</div>

{#each sheets as sheet (sheet.page_number)}
  {@const groups = groupEntriesByClient(sheet.entries)}
  <div class="sheet-page">
    <div class="mb-3 text-sm">
      <div class="flex flex-wrap gap-4">
        <span class="font-semibold">Zona: {sheet.zone_name}</span>
        <span>Fecha creación: {formatDate(sheet.created_at)}</span>
        <span>Página {sheet.page_number} de {sheets.length}</span>
      </div>
    </div>

    <table class="w-full border-collapse text-sm">
      <thead>
        <tr>
          <th
            rowspan="2"
            class="border border-border px-2 py-2 text-left align-top"
          >
            Cliente
          </th>
          {#each dayColumns as day}
            <th colspan="2" class="border border-border px-2 py-2 text-center">
              Día {day}
            </th>
          {/each}
        </tr>
        <tr>
          {#each dayColumns as _}
            <th class="border border-border px-1 py-2 text-center text-xs">C</th
            >
            <th class="border border-border px-1 py-2 text-center text-xs">F</th
            >
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each groups as group (group.clientId)}
          <tr>
            <td class="border border-border px-2 py-2 align-top font-semibold">
              {group.clientName} ({group.clientBalance})
            </td>
            {#each dayColumns as _}
              <td class="day-cell border border-border px-1 py-2"></td>
              <td class="day-cell border border-border px-1 py-2"></td>
            {/each}
          </tr>
          <tr>
            <td
              class="border border-border px-2 py-2 text-xs text-muted-foreground"
            >
              {group.clientAddress || "—"}
            </td>
            {#each dayColumns as _}
              <td class="day-cell border border-border px-1 py-2"></td>
              <td class="day-cell border border-border px-1 py-2"></td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/each}
