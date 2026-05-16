<script lang="ts">
  import * as Select from "$lib/components/ui/select/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Popover from "$lib/components/ui/popover/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Calendar } from "$lib/components/ui/calendar/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { toast } from "svelte-sonner";
  import { parseDate, type DateValue } from "@internationalized/date";
  import { GetZoneReport } from "wailsjs/go/service/ReportService";
  import { getZones, fetchZones } from "$lib/stores/zones.svelte.js";
  import type { dto } from "wailsjs/go/models";
  function formatDateToDDMMYYYY(isoDate: string): string {
    const [year, month, day] = isoDate.split("-");
    return `${day}-${month}-${year}`;
  }

  function parseFromDDMMYYYY(formatted: string): string {
    const [day, month, year] = formatted.split("-");
    return `${year}-${month}-${day}`;
  }

  function defaultEndDate() {
    return formatDateToDDMMYYYY(new Date().toISOString().slice(0, 10));
  }

  function defaultStartDate() {
    const d = new Date();
    d.setDate(d.getDate() - 30);
    return formatDateToDDMMYYYY(d.toISOString().slice(0, 10));
  }

  let zones = $derived(getZones());
  let selectedZoneId = $state<string | undefined>(undefined);
  let startDate = $state(defaultStartDate());
  let endDate = $state(defaultEndDate());
  let startDateValue = $derived.by(
    () => parseDate(parseFromDDMMYYYY(startDate)) as DateValue,
  );
  let endDateValue = $derived.by(
    () => parseDate(parseFromDDMMYYYY(endDate)) as DateValue,
  );
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
      const isoStart = parseFromDDMMYYYY(startDate);
      const isoEnd = parseFromDDMMYYYY(endDate);
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
      <Label>Zona</Label>
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

    <div class="grid gap-1">
      <Label for="report-start">Desde</Label>
      <Popover.Root>
        <Popover.Trigger>
          <Button
            variant="outline"
            class="w-40 justify-start text-left font-normal"
          >
            {startDate}
          </Button>
        </Popover.Trigger>
        <Popover.Content class="w-auto p-0" align="start">
          <Calendar
            type="single"
            value={startDateValue}
            onValueChange={(date: DateValue | undefined) => {
              if (date) {
                const stringValue = date.toString();
                const [year, month, day] = stringValue.split("-");
                startDate = `${day}-${month}-${year}`;
              }
            }}
          />
        </Popover.Content>
      </Popover.Root>
    </div>

    <div class="grid gap-1">
      <Label for="report-end">Hasta</Label>
      <Popover.Root>
        <Popover.Trigger>
          <Button
            variant="outline"
            class="w-40 justify-start text-left font-normal"
          >
            {endDate}
          </Button>
        </Popover.Trigger>
        <Popover.Content class="w-auto p-0" align="start">
          <Calendar
            type="single"
            value={endDateValue}
            onValueChange={(date: DateValue | undefined) => {
              if (date) {
                const stringValue = date.toString();
                const [year, month, day] = stringValue.split("-");
                endDate = `${day}-${month}-${year}`;
              }
            }}
          />
        </Popover.Content>
      </Popover.Root>
    </div>

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
            <Table.Cell class="text-right">{row.delivery_date}</Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  {/if}
</div>
