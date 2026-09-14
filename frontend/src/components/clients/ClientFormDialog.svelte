<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import * as Select from "$lib/components/ui/select/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { createClient, updateClient } from "$lib/stores/clients.svelte.js";
  import { toast } from "svelte-sonner";
  import type { model } from "wailsjs/go/models";

  let {
    open = $bindable(false),
    client,
    zones,
  }: { open: boolean; client?: model.Client; zones: model.Zone[] } = $props();

  let name = $state("");
  let address = $state("");
  let zoneId = $state<string | undefined>(undefined);
  let submitting = $state(false);

  let isEdit = $derived(!!client);
  let title = $derived(isEdit ? "Editar Cliente" : "Nuevo Cliente");

  $effect(() => {
    if (open) {
      name = client?.name ?? "";
      address = client?.address ?? "";
      zoneId = client?.zone_id ? client.zone_id : undefined;
    }
  });

  let zoneMap = $derived(new Map(zones.map((z) => [z.reference_uuid, z.name])));

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    const trimmedName = name.trim();
    if (!trimmedName || !zoneId) return;

    submitting = true;
    try {
      if (isEdit && client) {
        await updateClient(client.id, trimmedName, address.trim(), zoneId);
        toast.success("Cliente actualizado");
      } else {
        await createClient(trimmedName, address.trim(), zoneId);
        toast.success("Cliente creado");
      }
      open = false;
    } catch (err) {
      toast.error(String(err));
    } finally {
      submitting = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="sm:max-w-md">
    <Dialog.Header>
      <Dialog.Title>{title}</Dialog.Title>
    </Dialog.Header>
    <form onsubmit={handleSubmit} class="space-y-4">
      <div class="space-y-2">
        <Label for="client-name">Nombre</Label>
        <Input id="client-name" bind:value={name} placeholder="Nombre del cliente" required />
      </div>
      <div class="space-y-2">
        <Label for="client-address">Dirección</Label>
        <Input id="client-address" bind:value={address} placeholder="Dirección (opcional)" />
      </div>
      <div class="space-y-2">
        <Label>Zona</Label>
        <Select.Root type="single" bind:value={zoneId}>
          <Select.Trigger class="w-105 truncate">
            {#if zoneId}
              {zoneMap.get(zoneId) ?? "Seleccionar zona"}
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
      <Dialog.Footer>
        <Button variant="outline" type="button" onclick={() => (open = false)}>Cancelar</Button>
        <Button type="submit" disabled={submitting || !name.trim() || !zoneId}>
          {submitting ? "Guardando..." : "Guardar"}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
