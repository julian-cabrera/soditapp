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
    getPreviewMode,
    setPreviewMode,
    groupEntriesByClient,
    getWeekStart,
    getWeekEnd,
  } from "$lib/stores/delivery.svelte.js";
  import DeliveryConfirmDialog from "./DeliveryConfirmDialog.svelte";
  import DeliverySheetPreview from "./DeliverySheetPreview.svelte";
  import WeekSelector from "./WeekSelector.svelte";
  import type { dto } from "wailsjs/go/models";

  let selectedZoneId = $state<string | undefined>(undefined);
  let confirmOpen = $state(false);

  let zones = $derived(getZones());
  let sheets = $derived(getSheets());
  let loading = $derived(getDeliveryLoading());
  let currentPageIndex = $derived(getCurrentPage());
  let entryEdits = $derived(getEntryEdits());
  let currentSheet = $derived(sheets[currentPageIndex] as dto.DeliverySheetDTO | undefined);
  let previewMode = $derived(getPreviewMode());
  let weekStart = $derived(getWeekStart());
  let weekEnd = $derived(getWeekEnd());

  let weekLabel = $derived(
    `Semana del ${formatDate(weekStart)} al ${formatDate(weekEnd)}`,
  );

  let adjustmentCount = $derived(() => {
    let count = 0;
    for (const edit of entryEdits.values()) {
      if (edit.credit > 0) count++;
    }
    return count;
  });

  let groupedEntries = $derived(() => {
    if (!currentSheet) return [];
    return groupEntriesByClient(currentSheet.entries);
  });

  function formatDate(date: Date): string {
    const d = String(date.getDate()).padStart(2, "0");
    const m = String(date.getMonth() + 1).padStart(2, "0");
    const y = date.getFullYear();
    return `${d}/${m}/${y}`;
  }

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
</script>

{#if previewMode}
  <div class="space-y-4">
    <div class="print-hidden">
      <WeekSelector />
    </div>
    <DeliverySheetPreview
      {sheets}
      {weekLabel}
      onBack={() => setPreviewMode(false)}
    />
  </div>
{:else}
  <div class="space-y-4">
    <h2 class="text-lg font-semibold">Entregas</h2>

    <WeekSelector />

    <div class="flex items-center gap-2">
      <span class="text-sm text-muted-foreground">Zona:</span>
      <Select.Root type="single" bind:value={selectedZoneId}>
        <Select.Trigger class="w-48">
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

    {#if !selectedZoneId}
      <p class="text-sm text-muted-foreground">
        Seleccioná una zona para ver los clientes con saldo pendiente.
      </p>
    {:else if loading && sheets.length === 0}
      <p class="text-sm text-muted-foreground">Cargando...</p>
    {:else if sheets.length === 0}
      <p class="text-sm text-muted-foreground">
        No hay clientes con saldo pendiente en esta zona.
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
                    value={entryEdits.get(`${entry.client_id}-${entry.product_id}`)?.charged ?? 0}
                    oninput={(e) => handleInput(entry.client_id, entry.product_id, "charged", e)}
                  />
                </Table.Cell>
                <Table.Cell>
                  <Input
                    type="number"
                    min="0"
                    class="h-8 w-20"
                    value={entryEdits.get(`${entry.client_id}-${entry.product_id}`)?.credit ?? 0}
                    oninput={(e) => handleInput(entry.client_id, entry.product_id, "credit", e)}
                  />
                </Table.Cell>
              </Table.Row>
            {/each}
          {/each}
        </Table.Body>
      </Table.Root>

      <div class="flex justify-end gap-2">
        <Button
          variant="outline"
          disabled={sheets.length === 0}
          onclick={() => setPreviewMode(true)}
        >
          Ver Hoja
        </Button>
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

  <DeliveryConfirmDialog bind:open={confirmOpen} adjustmentCount={adjustmentCount()} />
{/if}
