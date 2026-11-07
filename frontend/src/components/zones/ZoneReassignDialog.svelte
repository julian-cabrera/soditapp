<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { Checkbox } from "$lib/components/ui/checkbox/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { toast } from "svelte-sonner";
  import { ListClientsByZone, MoveClientsToZone } from "wailsjs/go/service/ClientService";
  import { fetchZones } from "$lib/stores/zones.svelte.js";
  import { fetchClients } from "$lib/stores/clients.svelte.js";
  import type { model } from "wailsjs/go/models";

  let { open = $bindable(false), zones }: { open: boolean; zones: model.Zone[] } = $props();

  let step = $state<1 | 2>(1);
  let sourceZoneId = $state<string | undefined>(undefined);
  let targetZoneId = $state<string | undefined>(undefined);
  let clients = $state<model.Client[]>([]);
  let selectedIds = $state<Set<string>>(new Set());
  let loading = $state(false);
  let submitting = $state(false);

  let availableTargetZones = $derived(
    zones.filter((z) => z.reference_uuid !== sourceZoneId)
  );

  $effect(() => {
    if (!open) {
      step = 1;
      sourceZoneId = undefined;
      targetZoneId = undefined;
      clients = [];
      selectedIds = new Set();
    }
  });

  async function handleNext() {
    if (!sourceZoneId || !targetZoneId) return;
    loading = true;
    try {
      clients = await ListClientsByZone(sourceZoneId);
      selectedIds = new Set(clients.map((c) => c.reference_uuid));
      step = 2;
    } catch (e) {
      toast.error(String(e));
    } finally {
      loading = false;
    }
  }

  function toggleClient(refUUID: string) {
    const next = new Set(selectedIds);
    if (next.has(refUUID)) {
      next.delete(refUUID);
    } else {
      next.add(refUUID);
    }
    selectedIds = next;
  }

  async function handleConfirm() {
    if (selectedIds.size === 0 || !targetZoneId) return;
    submitting = true;
    try {
      await MoveClientsToZone(Array.from(selectedIds), targetZoneId);
      toast.success("Clientes reasignados correctamente");
      await Promise.all([fetchZones(), fetchClients()]);
      open = false;
    } catch (e) {
      toast.error(String(e));
    } finally {
      submitting = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="max-w-md">
    <Dialog.Header>
      <Dialog.Title>Reasignar Clientes</Dialog.Title>
      <Dialog.Description>
        {#if step === 1}
          Seleccioná la zona de origen y la zona de destino.
        {:else}
          Seleccioná los clientes a mover.
        {/if}
      </Dialog.Description>
    </Dialog.Header>

    {#if step === 1}
      <div class="space-y-4 py-2">
        <div class="space-y-1">
          <Label>Zona de origen</Label>
          <Select.Root type="single" bind:value={sourceZoneId}>
            <Select.Trigger class="w-full">
              {#if sourceZoneId}
                {zones.find((z) => z.reference_uuid === sourceZoneId)?.name ?? "Zona"}
              {:else}
                Seleccionar zona
              {/if}
            </Select.Trigger>
            <Select.Content>
              {#each zones as zone (zone.id)}
                <Select.Item value={zone.reference_uuid}>{zone.name}</Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </div>

        <div class="space-y-1">
          <Label>Zona de destino</Label>
          <Select.Root type="single" bind:value={targetZoneId}>
            <Select.Trigger class="w-full">
              {#if targetZoneId}
                {zones.find((z) => z.reference_uuid === targetZoneId)?.name ?? "Zona"}
              {:else}
                Seleccionar zona
              {/if}
            </Select.Trigger>
            <Select.Content>
              {#each availableTargetZones as zone (zone.id)}
                <Select.Item value={zone.reference_uuid}>{zone.name}</Select.Item>
              {/each}
            </Select.Content>
          </Select.Root>
        </div>
      </div>

      <Dialog.Footer>
        <Button variant="outline" onclick={() => (open = false)}>Cancelar</Button>
        <Button onclick={handleNext} disabled={!sourceZoneId || !targetZoneId || loading}>
          {loading ? "Cargando..." : "Siguiente"}
        </Button>
      </Dialog.Footer>

    {:else}
      <div class="py-2">
        {#if clients.length === 0}
          <p class="text-sm text-muted-foreground">No hay clientes activos en esta zona.</p>
        {:else}
          <ul class="max-h-60 overflow-y-auto space-y-2">
            {#each clients as client (client.id)}
              <li class="flex items-center gap-2">
                <Checkbox
                  id="client-{client.id}"
                  checked={selectedIds.has(client.reference_uuid)}
                  onCheckedChange={() => toggleClient(client.reference_uuid)}
                />
                <Label for="client-{client.id}" class="cursor-pointer font-normal">
                  {client.name}
                </Label>
              </li>
            {/each}
          </ul>
        {/if}
      </div>

      <Dialog.Footer>
        <Button variant="outline" onclick={() => (step = 1)}>Volver</Button>
        <Button
          onclick={handleConfirm}
          disabled={selectedIds.size === 0 || submitting}
        >
          {submitting ? "Moviendo..." : "Confirmar"}
        </Button>
      </Dialog.Footer>
    {/if}
  </Dialog.Content>
</Dialog.Root>
