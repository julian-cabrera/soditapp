<script lang="ts">
  import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
  import { unregisterClient } from "$lib/stores/clients.svelte.js";
  import { toast } from "svelte-sonner";
  import type { model } from "wailsjs/go/models";

  let { open = $bindable(false), client }: { open: boolean; client: model.Client } = $props();

  let processing = $state(false);

  async function handleUnregister() {
    processing = true;
    try {
      await unregisterClient(client.id);
      toast.success("Cliente dado de baja");
      open = false;
    } catch (err) {
      toast.error(String(err));
    } finally {
      processing = false;
    }
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Dar de baja</AlertDialog.Title>
      <AlertDialog.Description>
        ¿Estás seguro de que deseas dar de baja a <strong>{client.name}</strong>?
        El cliente no aparecerá en las listas activas.
      </AlertDialog.Description>
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel>Cancelar</AlertDialog.Cancel>
      <AlertDialog.Action
        onclick={handleUnregister}
        disabled={processing}
        class="bg-destructive text-destructive-foreground hover:bg-destructive/90"
      >
        {processing ? "Procesando..." : "Dar de baja"}
      </AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
