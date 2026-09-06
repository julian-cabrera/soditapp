<script lang="ts">
  import * as Select from "$lib/components/ui/select/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import DatePicker from "$components/ui/DatePicker.svelte";
  import { toast } from "svelte-sonner";
  import { GetZoneReport } from "wailsjs/go/service/ReportService";
  import { getZones, fetchZones } from "$lib/stores/zones.svelte.js";
  import { displayToISO, isoToDisplay, defaultStartDate, defaultEndDate } from "$lib/utils/date.js";
  import type { dto } from "wailsjs/go/models";

  let zones = $derived(getZones());
  let selectedZoneId = $state<string | undefined>(undefined);
  let startDate = $state(defaultStartDate());
  let endDate = $state(defaultEndDate());
  let rows = $state<dto.ZoneReportEntryDTO[]>([]);
  let loading = $state(false);

  let sortedRows = $derived(
    [...rows].sort((a, b) => {
      const dateCmp = a.delivery_date.localeCompare(b.delivery_date);
      if (dateCmp !== 0) return dateCmp;
      return a.client_name.localeCompare(b.client_name);
    }),
  );

  $effect(() => {
    fetchZones();
  });

  async function search() {
    if (!selectedZoneId) return;
    loading = true;
    try {
      const isoStart = displayToISO(startDate);
      const isoEnd = displayToISO(endDate);
      rows = await GetZoneReport(Number(selectedZoneId), isoStart, isoEnd);
    } catch (err) {
      toast.error(String(err));
    } finally {
      loading = false;
    }
  }
</script>

<div class="space-y-4">
  <h2 class="text-lg font-semibold">Reportes por Zona</h2>

  <div class="flex flex-wrap items-end gap-3">
    <div class="grid gap-1">
      <span class="text-sm font-medium">Zona</span>
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

    <DatePicker bind:value={startDate} label="Desde" />
    <DatePicker bind:value={endDate} label="Hasta" />

    <Button onclick={search} disabled={!selectedZoneId || loading}
      >Buscar</Button
    >
  </div>

  {#if loading}
    <p class="text-sm text-muted-foreground">Cargando...</p>
  {:else if sortedRows.length === 0 && selectedZoneId}
    <p class="text-sm text-muted-foreground">
      Sin resultados para el período seleccionado.
    </p>
  {:else if sortedRows.length > 0}
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Cliente</Table.Head>
          <Table.Head>Producto</Table.Head>
          <Table.Head class="text-right">Cobrado</Table.Head>
          <Table.Head class="text-right">Fiado</Table.Head>
          <Table.Head>Fecha</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#each sortedRows as row, i (i)}
          <Table.Row>
            <Table.Cell>{row.client_name}</Table.Cell>
            <Table.Cell>{row.product_name}</Table.Cell>
            <Table.Cell class="text-right">{row.charged_units}</Table.Cell>
            <Table.Cell class="text-right">{row.credit_units}</Table.Cell>
            <Table.Cell class="text-right">{isoToDisplay(row.delivery_date)}</Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  {/if}
</div>
