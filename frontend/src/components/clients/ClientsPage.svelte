<script lang="ts">
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { PencilEdit01Icon, Delete01Icon, PlusSignIcon, Coins01Icon, Download01Icon, TimeScheduleIcon } from "@hugeicons/core-free-icons";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import { toast } from "svelte-sonner";
  import { ExportClients, ExportBalances } from "wailsjs/go/service/ExportService";
  import {
    getClients,
    getClientsLoading,
    fetchClients,
    fetchClientsByZone,
  } from "$lib/stores/clients.svelte.js";
  import { getZones, fetchZones } from "$lib/stores/zones.svelte.js";
  import ClientFormDialog from "./ClientFormDialog.svelte";
  import ClientUnregisterDialog from "./ClientUnregisterDialog.svelte";
  import ClientBalancesDialog from "./ClientBalancesDialog.svelte";
  import ImportClientsDialog from "./ImportClientsDialog.svelte";
  import DeliveryHistoryDialog from "./DeliveryHistoryDialog.svelte";
  import type { model } from "wailsjs/go/models";

  async function exportClients() {
    try {
      await ExportClients();
      toast.success("Clientes exportados correctamente");
    } catch (e) {
      toast.error(String(e));
    }
  }

  async function exportBalances() {
    try {
      await ExportBalances();
      toast.success("Saldos exportados correctamente");
    } catch (e) {
      toast.error(String(e));
    }
  }

  let formOpen = $state(false);
  let unregisterOpen = $state(false);
  let balancesOpen = $state(false);
  let importOpen = $state(false);
  let historyOpen = $state(false);
  let selectedClient = $state<model.Client | undefined>(undefined);
  let balancesClient = $state<model.Client | undefined>(undefined);
  let historyClient = $state<model.Client | undefined>(undefined);
  let selectedZoneId = $state<string | undefined>(undefined);
  let searchQuery = $state("");

  let clients = $derived(getClients());
  let zones = $derived(getZones());
  let loading = $derived(getClientsLoading());
  let zoneMap = $derived(new Map(zones.map((z) => [z.id, z.name])));
  let filteredClients = $derived(
    searchQuery
      ? clients.filter((c) => {
          const q = searchQuery.toLowerCase();
          return c.name.toLowerCase().includes(q) || (c.address ?? "").toLowerCase().includes(q);
        })
      : clients
  );

  $effect(() => {
    fetchZones();
    fetchClients();
  });

  $effect(() => {
    if (selectedZoneId && selectedZoneId !== "all") {
      fetchClientsByZone(Number(selectedZoneId));
    } else {
      fetchClients();
    }
  });

  function openCreate() {
    selectedClient = undefined;
    formOpen = true;
  }

  function openEdit(client: model.Client) {
    selectedClient = client;
    formOpen = true;
  }

  function openUnregister(client: model.Client) {
    selectedClient = client;
    unregisterOpen = true;
  }

  function openBalances(client: model.Client) {
    balancesClient = client;
    balancesOpen = true;
  }

  function openHistory(client: model.Client) {
    historyClient = client;
    historyOpen = true;
  }
</script>

<div class="space-y-4">
  <div class="flex items-center justify-between">
    <h2 class="text-lg font-semibold">Clientes</h2>
    <div class="flex gap-2">
      <Button variant="outline" size="default" onclick={exportClients}>
        <HugeiconsIcon icon={Download01Icon} size={14} />
        Exportar Clientes
      </Button>
      <Button variant="outline" size="default" onclick={exportBalances}>
        <HugeiconsIcon icon={Download01Icon} size={14} />
        Exportar Saldos
      </Button>
      <Button variant="outline" size="default" onclick={() => (importOpen = true)}>
        Importar Clientes
      </Button>
      <Button size="default" onclick={openCreate}>
        <HugeiconsIcon icon={PlusSignIcon} size={14} />
        Nuevo Cliente
      </Button>
    </div>
  </div>

  <div class="flex items-center gap-2">
    <Input placeholder="Buscar por nombre o dirección..." bind:value={searchQuery} class="max-w-xs" />
    <span class="text-sm text-muted-foreground">Filtrar por zona:</span>
    <Select.Root type="single" bind:value={selectedZoneId}>
      <Select.Trigger class="w-48">
        {#if selectedZoneId && selectedZoneId !== "all"}
          {zoneMap.get(Number(selectedZoneId)) ?? "Zona"}
        {:else}
          Todas las zonas
        {/if}
      </Select.Trigger>
      <Select.Content>
        <Select.Item value="all">Todas las zonas</Select.Item>
        {#each zones as zone (zone.id)}
          <Select.Item value={String(zone.id)}>{zone.name}</Select.Item>
        {/each}
      </Select.Content>
    </Select.Root>
  </div>

  {#if loading && clients.length === 0}
    <p class="text-sm text-muted-foreground">Cargando...</p>
  {:else if filteredClients.length === 0}
    <p class="text-sm text-muted-foreground">No hay clientes registrados.</p>
  {:else}
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Nombre</Table.Head>
          <Table.Head>Dirección</Table.Head>
          <Table.Head>Zona</Table.Head>
          <Table.Head class="w-40 text-right">Acciones</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#each filteredClients as client (client.id)}
          <Table.Row>
            <Table.Cell>{client.name}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{client.address || "—"}</Table.Cell>
            <Table.Cell>
              <Badge variant="secondary">{zoneMap.get(client.zone_id) ?? "—"}</Badge>
            </Table.Cell>
            <Table.Cell class="text-right">
              <div class="flex justify-end gap-1">
                <Button variant="ghost" size="icon-xs" onclick={() => openHistory(client)} title="Ver Historial">
                  <HugeiconsIcon icon={TimeScheduleIcon} size={14} />
                </Button>
                <Button variant="ghost" size="icon-xs" onclick={() => openBalances(client)} title="Saldos">
                  <HugeiconsIcon icon={Coins01Icon} size={14} />
                </Button>
                <Button variant="ghost" size="icon-xs" onclick={() => openEdit(client)} title="Editar">
                  <HugeiconsIcon icon={PencilEdit01Icon} size={14} />
                </Button>
                <Button variant="ghost" size="icon-xs" onclick={() => openUnregister(client)} title="Dar de baja">
                  <HugeiconsIcon icon={Delete01Icon} size={14} />
                </Button>
              </div>
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  {/if}
</div>

<ClientFormDialog bind:open={formOpen} client={selectedClient} {zones} />

{#if selectedClient}
  <ClientUnregisterDialog bind:open={unregisterOpen} client={selectedClient} />
{/if}

{#if balancesClient}
  <ClientBalancesDialog bind:open={balancesOpen} client={balancesClient} />
{/if}

<ImportClientsDialog bind:open={importOpen} />

{#if historyClient}
  <DeliveryHistoryDialog bind:open={historyOpen} clientId={historyClient.id} clientName={historyClient.name} />
{/if}
