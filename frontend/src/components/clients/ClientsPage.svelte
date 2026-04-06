<script lang="ts">
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { PencilEdit01Icon, Delete01Icon, PlusSignIcon, Coins01Icon } from "@hugeicons/core-free-icons";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
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
  import type { model } from "wailsjs/go/models";

  let formOpen = $state(false);
  let unregisterOpen = $state(false);
  let balancesOpen = $state(false);
  let selectedClient = $state<model.Client | undefined>(undefined);
  let balancesClient = $state<model.Client | undefined>(undefined);
  let selectedZoneId = $state<string | undefined>(undefined);

  let clients = $derived(getClients());
  let zones = $derived(getZones());
  let loading = $derived(getClientsLoading());
  let zoneMap = $derived(new Map(zones.map((z) => [z.id, z.name])));

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
</script>

<div class="space-y-4">
  <div class="flex items-center justify-between">
    <h2 class="text-lg font-semibold">Clientes</h2>
    <Button size="default" onclick={openCreate}>
      <HugeiconsIcon icon={PlusSignIcon} size={14} />
      Nuevo Cliente
    </Button>
  </div>

  <div class="flex items-center gap-2">
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
  {:else if clients.length === 0}
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
        {#each clients as client (client.id)}
          <Table.Row>
            <Table.Cell>{client.name}</Table.Cell>
            <Table.Cell class="text-muted-foreground">{client.address || "—"}</Table.Cell>
            <Table.Cell>
              <Badge variant="secondary">{zoneMap.get(client.zone_id) ?? "—"}</Badge>
            </Table.Cell>
            <Table.Cell class="text-right">
              <div class="flex justify-end gap-1">
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
