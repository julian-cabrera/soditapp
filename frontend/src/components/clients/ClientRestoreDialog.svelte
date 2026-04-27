<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { RotateLeftIcon } from "@hugeicons/core-free-icons";
  import { toast } from "svelte-sonner";
  import {
    restoreClient,
    listInactiveClients,
  } from "$lib/stores/clients.svelte.js";
  import { getZones } from "$lib/stores/zones.svelte.js";
  import type { model } from "wailsjs/go/models";

  let { open = $bindable(false) }: { open: boolean } = $props();

  let inactiveClients = $state<model.Client[]>([]);
  let loading = $state(false);
  let restoring = $state(false);
  let zones = $derived(getZones());
  let zoneMap = $derived(new Map(zones.map((z) => [z.id, z.name])));

  $effect(() => {
    if (open) {
      loadInactiveClients();
    }
  });

  async function loadInactiveClients() {
    loading = true;
    try {
      inactiveClients = await listInactiveClients();
    } catch (err) {
      toast.error(String(err));
    } finally {
      loading = false;
    }
  }

  async function handleRestore(clientId: number) {
    restoring = true;
    try {
      await restoreClient(clientId);
      toast.success("Cliente restaurado");
      inactiveClients = await listInactiveClients();
    } catch (err) {
      toast.error(String(err));
    } finally {
      restoring = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="sm:max-w-2xl">
    <Dialog.Header>
      <Dialog.Title>Restaurar Cliente</Dialog.Title>
      <Dialog.Description>
        Selecciona un cliente dado de baja para reactivarlo
      </Dialog.Description>
    </Dialog.Header>
    {#if loading}
      <p class="text-sm text-muted-foreground">Cargando clientes...</p>
    {:else if inactiveClients.length === 0}
      <div
        class="flex flex-col items-center gap-2 rounded-md border border-dashed p-6 text-sm text-muted-foreground"
      >
        <p>No hay clientes dados de baja para restaurar.</p>
      </div>
    {:else}
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.Head>Nombre</Table.Head>
            <Table.Head>Dirección</Table.Head>
            <Table.Head>Zona</Table.Head>
            <Table.Head class="w-24 text-right">Acción</Table.Head>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {#each inactiveClients as client (client.id)}
            <Table.Row>
              <Table.Cell>{client.name}</Table.Cell>
              <Table.Cell class="text-muted-foreground"
                >{client.address || "—"}</Table.Cell
              >
              <Table.Cell>
                <Badge variant="secondary"
                  >{zoneMap.get(client.zone_id) ?? "—"}</Badge
                >
              </Table.Cell>
              <Table.Cell class="text-right">
                <Button
                  variant="ghost"
                  size="icon"
                  disabled={restoring}
                  onclick={() => handleRestore(client.id)}
                  title="Restaurar cliente"
                >
                  <HugeiconsIcon icon={RotateLeftIcon} size={14} />
                </Button>
              </Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    {/if}
    <Dialog.Footer>
      <Button variant="outline" type="button" onclick={() => (open = false)}
        >Cerrar</Button
      >
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
