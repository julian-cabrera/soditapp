<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import DatePicker from "$components/ui/DatePicker.svelte";
  import { toast } from "svelte-sonner";
  import { GetClientHistory } from "wailsjs/go/service/ReportService";
  import { displayToISO, isoToDisplay, defaultStartDate, defaultEndDate } from "$lib/utils/date.js";
  import type { dto } from "wailsjs/go/models";

  let {
    open = $bindable(false),
    clientId,
    clientName,
  }: { open: boolean; clientId: number; clientName: string } = $props();

  let startDate = $state(defaultStartDate());
  let endDate = $state(defaultEndDate());
  let rows = $state<dto.DeliveryHistoryEntryDTO[]>([]);
  let loading = $state(false);

  $effect(() => {
    if (open) {
      startDate = defaultStartDate();
      endDate = defaultEndDate();
      rows = [];
    }
  });

  async function search() {
    loading = true;
    try {
      rows = await GetClientHistory(clientId, displayToISO(startDate), displayToISO(endDate));
    } catch (err) {
      toast.error(String(err));
    } finally {
      loading = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="sm:max-w-2xl">
    <Dialog.Header>
      <Dialog.Title>Historial de entregas — {clientName}</Dialog.Title>
    </Dialog.Header>

    <div class="flex items-end gap-3">
      <DatePicker bind:value={startDate} label="Desde" />
      <DatePicker bind:value={endDate} label="Hasta" />
      <Button onclick={search} disabled={loading}>Buscar</Button>
    </div>

    {#if loading}
      <p class="text-sm text-muted-foreground">Cargando...</p>
    {:else if rows.length === 0}
      <p class="text-sm text-muted-foreground">Sin resultados para el período seleccionado.</p>
    {:else}
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.Head>Fecha</Table.Head>
            <Table.Head>Producto</Table.Head>
            <Table.Head class="text-right">Cobrado</Table.Head>
            <Table.Head class="text-right">Fiado</Table.Head>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {#each rows as row, i (i)}
            <Table.Row>
              <Table.Cell>{isoToDisplay(row.delivery_date)}</Table.Cell>
              <Table.Cell>{row.product_name}</Table.Cell>
              <Table.Cell class="text-right">{row.charged_units}</Table.Cell>
              <Table.Cell class="text-right">{row.credit_units}</Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    {/if}

    <Dialog.Footer>
      <Button variant="outline" onclick={() => (open = false)}>Cerrar</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
