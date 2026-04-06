<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { toast } from "svelte-sonner";
  import { GetClientHistory } from "wailsjs/go/service/ReportService";
  import type { dto } from "wailsjs/go/models";

  let {
    open = $bindable(false),
    clientId,
    clientName,
  }: { open: boolean; clientId: number; clientName: string } = $props();

  function defaultEndDate() {
    return new Date().toISOString().slice(0, 10);
  }

  function defaultStartDate() {
    const d = new Date();
    d.setDate(d.getDate() - 30);
    return d.toISOString().slice(0, 10);
  }

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
      rows = await GetClientHistory(clientId, startDate, endDate);
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
      <div class="grid gap-1">
        <Label for="start-date">Desde</Label>
        <Input id="start-date" type="date" bind:value={startDate} class="w-40" />
      </div>
      <div class="grid gap-1">
        <Label for="end-date">Hasta</Label>
        <Input id="end-date" type="date" bind:value={endDate} class="w-40" />
      </div>
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
              <Table.Cell>{row.delivery_date}</Table.Cell>
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
