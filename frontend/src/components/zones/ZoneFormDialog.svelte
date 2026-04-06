<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { createZone, updateZone } from "$lib/stores/zones.svelte.js";
  import { toast } from "svelte-sonner";
  import type { model } from "wailsjs/go/models";

  let { open = $bindable(false), zone }: { open: boolean; zone?: model.Zone } = $props();

  let name = $state("");
  let submitting = $state(false);

  let isEdit = $derived(!!zone);
  let title = $derived(isEdit ? "Renombrar Zona" : "Nueva Zona");

  $effect(() => {
    if (open) {
      name = zone?.name ?? "";
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    const trimmed = name.trim();
    if (!trimmed) return;

    submitting = true;
    try {
      if (isEdit && zone) {
        await updateZone(zone.id, trimmed);
        toast.success("Zona actualizada");
      } else {
        await createZone(trimmed);
        toast.success("Zona creada");
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
        <Label for="zone-name">Nombre</Label>
        <Input id="zone-name" bind:value={name} placeholder="Nombre de la zona" required />
      </div>
      <Dialog.Footer>
        <Button variant="outline" type="button" onclick={() => (open = false)}>Cancelar</Button>
        <Button type="submit" disabled={submitting || !name.trim()}>
          {submitting ? "Guardando..." : "Guardar"}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
