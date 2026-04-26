<script lang="ts">
  import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
  import { deleteZone } from "$lib/stores/zones.svelte.js";
  import { toast } from "svelte-sonner";
  import type { model } from "wailsjs/go/models";

  let { open = $bindable(false), zone }: { open: boolean; zone: model.Zone } = $props();

  let deleting = $state(false);

  async function handleDelete() {
    deleting = true;
    try {
      await deleteZone(zone.id);
      toast.success("Zona eliminada");
      open = false;
    } catch (err) {
      toast.error(String(err));
    } finally {
      deleting = false;
    }
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Content class="overflow-auto">
    <AlertDialog.Header>
      <AlertDialog.Title>Eliminar zona</AlertDialog.Title>
      <AlertDialog.Description>
        ¿Estás seguro de que deseas eliminar la zona <strong>{zone.name}</strong>?
        <br />
        Esta acción no se puede deshacer.
      </AlertDialog.Description>
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel>Cancelar</AlertDialog.Cancel>
      <AlertDialog.Action
        onclick={handleDelete}
        disabled={deleting}
        class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
      >
        {deleting ? "Eliminando..." : "Eliminar"}
      </AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
