<script lang="ts">
  import { Button } from "$lib/components/ui/button/index.js";
  import { groupEntriesByClient } from "$lib/stores/delivery.svelte.js";
  import type { dto } from "wailsjs/go/models";

  let {
    sheets,
    dateLabel,
    onBack,
  }: {
    sheets: dto.DeliverySheetDTO[];
    dateLabel: string;
    onBack: () => void;
  } = $props();

  const dayColumns = [1, 2, 3, 4];
</script>

<div class="print-hidden mb-4 flex items-center justify-between">
  <Button variant="outline" onclick={onBack}>Volver</Button>
  <Button onclick={() => window.print()}>Imprimir</Button>
</div>

{#each sheets as sheet (sheet.page_number)}
  {@const groups = groupEntriesByClient(sheet.entries)}
  <div class="sheet-page">
    <div class="mb-2 flex items-center justify-between text-sm">
      <span class="font-semibold">Zona: {sheet.zone_name}</span>
      <span>{dateLabel}</span>
      <span>Página {sheet.page_number} de {sheets.length}</span>
    </div>

    <table class="w-full border-collapse text-sm">
      <thead>
        <tr>
          <th rowspan="2" class="w-[30%] border border-border px-2 py-1 text-left">
            Cliente
          </th>
          {#each dayColumns as day}
            <th colspan="2" class="border border-border px-2 py-1 text-center">
              Día {day}
            </th>
          {/each}
        </tr>
        <tr>
          {#each dayColumns as _}
            <th class="border border-border px-1 py-1 text-center text-xs">C</th>
            <th class="border border-border px-1 py-1 text-center text-xs">F</th>
          {/each}
        </tr>
      </thead>
      <tbody>
        {#each groups as group (group.clientId)}
          <tr>
            <td class="w-[30%] border border-border px-2 py-1 align-top">
              <div class="font-semibold">
                {group.clientName}
                <span class="ml-1 text-xs font-normal text-muted-foreground">
                  ({group.clientBalance})
                </span>
              </div>
              <div class="text-xs text-muted-foreground">
                {group.clientAddress || "—"}
              </div>
            </td>
            {#each dayColumns as _}
              <td class="day-cell border border-border px-1 py-1"></td>
              <td class="day-cell border border-border px-1 py-1"></td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
{/each}
