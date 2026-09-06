<script lang="ts">
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import {
    ArrowLeft01Icon,
    ArrowRight01Icon,
    Tick01Icon,
  } from "@hugeicons/core-free-icons";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import { getZones, fetchZones } from "$lib/stores/zones.svelte.js";
  import {
    getSheets,
    getDeliveryLoading,
    getCurrentPage,
    getEntryEdits,
    fetchDeliverySheets,
    setCurrentPage,
    updateEntryEdit,
    groupEntriesByClient,
    getWeekStartDate,
    getWeekEndDate,
    formatDateShort,
  } from "$lib/stores/delivery.svelte.js";
  import DeliveryConfirmDialog from "./DeliveryConfirmDialog.svelte";
  import WeekSelector from "./WeekSelector.svelte";
  import type { dto } from "wailsjs/go/models";

  let selectedZoneId = $state<string | undefined>(undefined);
  let confirmOpen = $state(false);
  let activeSection = $state<"load" | "print">("load");

  let zones = $derived(getZones());
  let sheets = $derived(getSheets());
  let loading = $derived(getDeliveryLoading());
  let currentPageIndex = $derived(getCurrentPage());
  let entryEdits = $derived(getEntryEdits());
  let currentSheet = $derived(
    sheets[currentPageIndex] as dto.DeliverySheetDTO | undefined,
  );
  let weekStart = $derived(formatDateShort(getWeekStartDate()));
  let weekEnd = $derived(formatDateShort(getWeekEndDate()));

  let adjustmentCount = $derived(() => {
    let count = 0;
    for (const edit of entryEdits.values()) {
      if (edit.charged > 0 || edit.credit > 0) count++;
    }
    return count;
  });

  let groupedEntries = $derived(() => {
    if (!currentSheet) return [];
    return groupEntriesByClient(currentSheet.entries);
  });

  $effect(() => {
    fetchZones();
  });

  $effect(() => {
    if (selectedZoneId) {
      fetchDeliverySheets(Number(selectedZoneId));
    }
  });

  function handleInput(
    clientId: number,
    productId: number,
    field: "charged" | "credit",
    e: Event,
  ) {
    const target = e.currentTarget as HTMLInputElement;
    updateEntryEdit(clientId, productId, field, Number(target.value) || 0);
  }

  function buildSheetHTML(
    sheet: dto.DeliverySheetDTO,
    sheetTotal: number,
    startIndex: number,
  ): string {
    const groups = groupEntriesByClient(sheet.entries);
    const dayColumns = [1, 2, 3, 4];
    const dayHeaders = dayColumns
      .map(() => `<th colspan="2" class="sheet-th">‎</th>`)
      .join("");
    const subHeaders = dayColumns
      .map(
        () => `<th class="sheet-th-sub">C</th><th class="sheet-th-sub">F</th>`,
      )
      .join("");
    const dayCells = dayColumns
      .map(
        () =>
          `<td class="sheet-td sheet-td-day"></td><td class="sheet-td sheet-td-day"></td>`,
      )
      .join("");

    const rows = groups
      .map(
        (g, i) => `
        <tr>
          <td class="sheet-td sheet-td-name"><span class="sheet-td-idx">${startIndex + i + 1}</span>${g.clientName}<br/><span class="sheet-td-addr">${g.clientAddress || "—"}</span></td>
          <td class="sheet-td sheet-td-num">${g.clientBalance}</td>
          ${dayCells}
        </tr>`,
      )
      .join("");

    return `
      <div class="sheet-page">
        <div class="sheet-page-header">
          <div class="sheet-header-left">
            <div><b>Fecha inicio</b>: ${weekStart}</div>
            <div><b>Fecha fin</b>: .....................................................</div>
          </div>
          <div class="sheet-header-center font-semibold"><b>Zona</b>: ${sheet.zone_name}</div>
          <div class="sheet-header-right"><b>Página</b> ${sheet.page_number} de ${sheetTotal}</div>
        </div>
        <table class="sheet-table">
          <colgroup>
            <col style="width:35%">
            <col style="width:5%">
            <col style="width:7.5%"><col style="width:7.5%">
            <col style="width:7.5%"><col style="width:7.5%">
            <col style="width:7.5%"><col style="width:7.5%">
            <col style="width:7.5%"><col style="width:7.5%">
          </colgroup>
          <thead>
            <tr>
              <th colspan="2" class="sheet-th"></th>
              ${dayHeaders}
            </tr>
            <tr>
              <th class="sheet-th-sub"></th>
              <th class="sheet-th-sub"></th>
              ${subHeaders}
            </tr>
          </thead>
          <tbody>${rows}</tbody>
        </table>
      </div>`;
  }

  function printDeliverySheet() {
    if (sheets.length === 0) return;
    let offset = 0;
    const pages = sheets.map((s) => {
      const html = buildSheetHTML(s, sheets.length, offset);
      offset += groupEntriesByClient(s.entries).length;
      return html;
    }).join("");
    const html = `<!DOCTYPE html>
                  <html lang="es">
                  <head>
                  <meta charset="UTF-8">
                  <title>Planilla de Entrega</title>
                  <style>
                    @page { size: 210mm 297mm; margin: 10mm; }
                    * { margin: 0; padding: 0; box-sizing: border-box; }
                    body { background: white; font-family: system-ui, -apple-system, sans-serif; }
                    .sheet-page {
                      width: 210mm; height: 297mm;
                      padding: 10mm;
                      page-break-after: always; break-after: page;
                      overflow: hidden;
                    }
                    .sheet-page:last-child { page-break-after: auto; break-after: auto; }
                    .sheet-page-header { display: flex; justify-content: space-between; align-items: flex-start; font-size: 11pt; margin-bottom: 0.5rem; }
                    .sheet-header-left { text-align: left; line-height: 1.6; }
                    .sheet-header-center { text-align: center; flex: 1; }
                    .sheet-header-right { text-align: right; }
                    .sheet-table { width: 100%; border-collapse: collapse; table-layout: fixed; font-size: 10pt; }
                    .sheet-table tbody tr { height: 10mm; }
                    .sheet-th { border: 1px solid #000; padding: 3px 6px; background: #eee; text-align: center; vertical-align: middle; }
                    .sheet-th-sub { border: 1px solid #000; padding: 2px 4px; background: #eee; text-align: center; font-size: 9pt; }
                    .sheet-td { border: 1px solid #000; padding: 3px 6px; vertical-align: top; }
                    .sheet-td-name { font-weight: 600; width: 35%; }
                    .sheet-td-idx { display: inline-block; width: 20px; text-align: right; margin-right: 4px; font-weight: 400; color: #666; }
                    .sheet-td-addr { font-weight: 400; font-size: 9pt; color: #666; }
                    .sheet-td-num { width: 5%; text-align: right; vertical-align: top; font-size: 10pt; }
                    @media print { .sheet-page { padding: 0; } }
                  </style>
                  </head>
                  <body>
                  <div>${pages}</div>
                  </body>
                  </html>`;

    const printWindow = window.open("", "_blank", "width=800,height=600");
    if (printWindow) {
      printWindow.document.write(html);
      printWindow.document.close();
      printWindow.focus();
      printWindow.print();
    }
  }
</script>

<div class="space-y-4">
  <h2 class="text-lg font-semibold">Entregas</h2>

  <div class="flex flex-wrap items-center gap-2">
    <Button
      variant={activeSection === "load" ? "default" : "outline"}
      size="sm"
      onclick={() => (activeSection = "load")}
    >
      Carga de planilla
    </Button>
    <Button
      variant={activeSection === "print" ? "default" : "outline"}
      size="sm"
      onclick={() => (activeSection = "print")}
    >
      Imprimir planilla
    </Button>
  </div>

  <div class="grid gap-4 md:grid-cols-[minmax(220px,_1fr)_minmax(220px,_1fr)]">
    <div class="space-y-2">
      <Select.Root type="single" bind:value={selectedZoneId}>
        <Select.Trigger class="w-full">
          {#if selectedZoneId}
            {zones.find((z) => String(z.id) === selectedZoneId)?.name ?? "Zona"}
          {:else}
            Seleccionar zona
          {/if}
        </Select.Trigger>
        <Select.Content>
          {#each zones as zone (zone.id)}
            <Select.Item value={String(zone.id)}>{zone.name}</Select.Item>
          {/each}
        </Select.Content>
      </Select.Root>
    </div>

    <div class="space-y-2">
      <WeekSelector />
    </div>
  </div>

  {#if activeSection === "load"}
    <div class="space-y-4">
      {#if !selectedZoneId}
        <p class="text-sm text-muted-foreground">
          Seleccioná una zona para ver los clientes.
        </p>
      {:else if loading && sheets.length === 0}
        <p class="text-sm text-muted-foreground">Cargando...</p>
      {:else if sheets.length === 0}
        <p class="text-sm text-muted-foreground">
          No hay clientes registrados en esta zona.
        </p>
      {:else}
        {#if sheets.length > 1}
          <div class="flex items-center gap-2">
            <Button
              variant="outline"
              size="sm"
              disabled={currentPageIndex === 0}
              onclick={() => setCurrentPage(currentPageIndex - 1)}
            >
              <HugeiconsIcon icon={ArrowLeft01Icon} size={14} />
            </Button>
            <span class="text-sm text-muted-foreground">
              Hoja {currentPageIndex + 1} de {sheets.length}
            </span>
            <Button
              variant="outline"
              size="sm"
              disabled={currentPageIndex === sheets.length - 1}
              onclick={() => setCurrentPage(currentPageIndex + 1)}
            >
              <HugeiconsIcon icon={ArrowRight01Icon} size={14} />
            </Button>
          </div>
        {/if}

        <Table.Root>
          <Table.Header>
            <Table.Row>
              <Table.Head>Cliente</Table.Head>
              <Table.Head>Dirección</Table.Head>
              <Table.Head class="text-right">Saldo</Table.Head>
              <Table.Head>Producto</Table.Head>
              <Table.Head class="w-24">Cobradas</Table.Head>
              <Table.Head class="w-24">Fiadas</Table.Head>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {#each groupedEntries() as group (group.clientId)}
              {#each group.entries as entry, i (entry.product_id)}
                <Table.Row class={i === 0 ? "border-t-2" : ""}>
                  {#if i === 0}
                    <Table.Cell class="font-medium align-top">
                      {group.clientName}
                    </Table.Cell>
                    <Table.Cell class="text-muted-foreground align-top">
                      {group.clientAddress || "—"}
                    </Table.Cell>
                    <Table.Cell class="text-right align-top">
                      <Badge variant="destructive">{group.clientBalance}</Badge>
                    </Table.Cell>
                  {:else}
                    <Table.Cell></Table.Cell>
                    <Table.Cell></Table.Cell>
                    <Table.Cell></Table.Cell>
                  {/if}
                  <Table.Cell>{entry.product_name}</Table.Cell>
                  <Table.Cell>
                    <Input
                      type="number"
                      min="0"
                      class="h-8 w-20"
                      value={entryEdits.get(
                        `${entry.client_id}-${entry.product_id}`,
                      )?.charged ?? 0}
                      oninput={(e) =>
                        handleInput(
                          entry.client_id,
                          entry.product_id,
                          "charged",
                          e,
                        )}
                    />
                  </Table.Cell>
                  <Table.Cell>
                    <Input
                      type="number"
                      min="0"
                      class="h-8 w-20"
                      value={entryEdits.get(
                        `${entry.client_id}-${entry.product_id}`,
                      )?.credit ?? 0}
                      oninput={(e) =>
                        handleInput(
                          entry.client_id,
                          entry.product_id,
                          "credit",
                          e,
                        )}
                    />
                  </Table.Cell>
                </Table.Row>
              {/each}
            {/each}
          </Table.Body>
        </Table.Root>

        <div class="flex justify-end gap-2">
          <Button
            disabled={adjustmentCount() === 0 || loading}
            onclick={() => (confirmOpen = true)}
          >
            <HugeiconsIcon icon={Tick01Icon} size={14} />
            Confirmar Entrega
          </Button>
        </div>
      {/if}
    </div>
  {:else}
    <div class="space-y-4">
      <p class="text-sm text-muted-foreground">
        Seleccioná una zona y fecha para generar la hoja de ruta.
      </p>
      <div class="flex flex-wrap items-center gap-2">
        <Button
          variant="default"
          disabled={!selectedZoneId || sheets.length === 0}
          onclick={printDeliverySheet}
        >
          Imprimir plantilla
        </Button>
        {#if selectedZoneId}
          <span class="text-sm text-muted-foreground">
            Zona seleccionada: {zones.find(
              (z) => String(z.id) === selectedZoneId,
            )?.name}
          </span>
        {/if}
      </div>
      {#if !selectedZoneId}
        <p class="text-sm text-muted-foreground">
          Debes seleccionar una zona antes de imprimir.
        </p>
      {/if}
    </div>
  {/if}

  <DeliveryConfirmDialog
    bind:open={confirmOpen}
    adjustmentCount={adjustmentCount()}
  />
</div>
