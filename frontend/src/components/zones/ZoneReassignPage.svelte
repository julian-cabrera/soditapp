<script lang="ts">
  import * as Select from "$lib/components/ui/select/index.js";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { toast } from "svelte-sonner";
  import { ListClientsByZone, MoveClientsToZone } from "wailsjs/go/service/ClientService";
  import { getZones, fetchZones } from "$lib/stores/zones.svelte.js";
  import { fetchClients } from "$lib/stores/clients.svelte.js";
  import type { model } from "wailsjs/go/models";

  let zones = $derived(getZones());

  let sourceZoneId = $state<string | undefined>(undefined);
  let targetZoneId = $state<string | undefined>(undefined);
  let clients = $state<model.Client[]>([]);
  let selectedIds = $state<Set<number>>(new Set());
  let loading = $state(false);
  let submitting = $state(false);

  let availableTargetZones = $derived(
    zones.filter((z) => String(z.id) !== sourceZoneId)
  );

  $effect(() => {
    fetchZones();
  });

  $effect(() => {
    if (sourceZoneId) {
      loadClients(sourceZoneId);
    } else {
      clients = [];
      selectedIds = new Set();
    }
  });

  async function loadClients(zoneId: string) {
    loading = true;
    targetZoneId = undefined;
    selectedIds = new Set();
    try {
      clients = await ListClientsByZone(Number(zoneId));
      selectedIds = new Set(clients.map((c) => c.id));
    } catch (e) {
      toast.error(String(e));
    } finally {
      loading = false;
    }
  }

  function toggleClient(id: number) {
    const next = new Set(selectedIds);
    if (next.has(id)) {
      next.delete(id);
    } else {
      next.add(id);
    }
    selectedIds = next;
  }

  function toggleAll() {
    if (selectedIds.size === clients.length) {
      selectedIds = new Set();
    } else {
      selectedIds = new Set(clients.map((c) => c.id));
    }
  }

  async function handleTransfer() {
    if (selectedIds.size === 0 || !targetZoneId) return;
    submitting = true;
    try {
      await MoveClientsToZone(Array.from(selectedIds), Number(targetZoneId));
      toast.success("Clientes reasignados correctamente");
      await Promise.all([fetchZones(), fetchClients()]);
      clients = await ListClientsByZone(Number(sourceZoneId));
      selectedIds = new Set(clients.map((c) => c.id));
      targetZoneId = undefined;
    } catch (e) {
      toast.error(String(e));
    } finally {
      submitting = false;
    }
  }
</script>

<div class="space-y-6">
  <h2 class="text-lg font-semibold">Reasignar Clientes</h2>

  <div class="max-w-xs space-y-1">
    <Label>Zona de origen</Label>
    <Select.Root type="single" bind:value={sourceZoneId}>
      <Select.Trigger class="w-full">
        {#if sourceZoneId}
          {zones.find((z) => String(z.id) === sourceZoneId)?.name ?? "Zona"}
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

  {#if sourceZoneId}
    {#if loading}
      <p class="text-sm text-muted-foreground">Cargando clientes...</p>
    {:else if clients.length === 0}
      <p class="text-sm text-muted-foreground">No hay clientes activos en esta zona.</p>
    {:else}
      <div class="space-y-3">
        <div class="flex items-center gap-2">
          <Checkbox
            id="select-all"
            checked={selectedIds.size === clients.length}
            onCheckedChange={toggleAll}
          />
          <Label for="select-all" class="cursor-pointer text-sm text-muted-foreground">
            Seleccionar todos ({clients.length})
          </Label>
        </div>

        <ul class="max-h-80 overflow-y-auto space-y-2 rounded-md border p-3">
          {#each clients as client (client.id)}
            <li class="flex items-center gap-2">
              <Checkbox
                id="client-{client.id}"
                checked={selectedIds.has(client.id)}
                onCheckedChange={() => toggleClient(client.id)}
              />
              <Label for="client-{client.id}" class="cursor-pointer font-normal">
                {client.name}
              </Label>
            </li>
          {/each}
        </ul>
      </div>

      <div class="flex items-center gap-3">
        <div class="w-48">
          <Select.Root type="single" bind:value={targetZoneId}>
            <Select.Trigger class="w-full">
              {#if targetZoneId}
                {zones.find((z) => String(z.id) === targetZoneId)?.name ?? "Zona"}
              {:else}
                Mover a...
              {/if}
            </Select.Trigger>
            <Select.Content>
              {#each availableTargetZones as zone (zone.id)}
                <Select.Item value={String(zone.id)}>{zone.name}</Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </div>

        <Button
          onclick={handleTransfer}
          disabled={selectedIds.size === 0 || !targetZoneId || submitting}
        >
          {submitting ? "Moviendo..." : `Transferir ${selectedIds.size > 0 ? `(${selectedIds.size})` : ""}`}
        </Button>
      </div>
    {/if}
  {/if}
</div>
