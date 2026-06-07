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
  import * as Dialog from "$lib/components/ui/dialog/index.js";
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
    getSelectedDate,
  } from "$lib/stores/delivery.svelte.js";
  import DeliveryConfirmDialog from "./DeliveryConfirmDialog.svelte";
  import DeliverySheetPreview from "./DeliverySheetPreview.svelte";
  import WeekSelector from "./WeekSelector.svelte";
  import type { dto } from "wailsjs/go/models";

  let selectedZoneId = $state<string | undefined>(undefined);
  let confirmOpen = $state(false);
  let printOpen = $state(false);
  let activeSection = $state<"load" | "print">("load");

  let zones = $derived(getZones());
  let sheets = $derived(getSheets());
  let loading = $derived(getDeliveryLoading());
  let currentPageIndex = $derived(getCurrentPage());
  let entryEdits = $derived(getEntryEdits());
  let currentSheet = $derived(
    sheets[currentPageIndex] as dto.DeliverySheetDTO | undefined,
  );
  let selectedDate = $derived(getSelectedDate());

  let dateLabel = $derived(formatDate(selectedDate));

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
      <span class="text-sm text-muted-foreground">Zona:</span>
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
      <span class="text-sm text-muted-foreground">Fecha:</span>
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
          disabled={!selectedZoneId}
          onclick={() => (printOpen = true)}
        >
          Imprimir planilla
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

  <Dialog.Root bind:open={printOpen}>
    <Dialog.Content class="max-w-5xl w-full">
      <DeliverySheetPreview {sheets} onBack={() => (printOpen = false)} />
    </Dialog.Content>
  </Dialog.Root>

  <DeliveryConfirmDialog
    bind:open={confirmOpen}
    adjustmentCount={adjustmentCount()}
  />
</div>
