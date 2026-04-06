<script lang="ts">
  import * as AlertDialog from "$lib/components/ui/alert-dialog/index.js";
  import { confirmDelivery } from "$lib/stores/delivery.svelte.js";
  import { toast } from "svelte-sonner";

  let {
    open = $bindable(false),
    adjustmentCount,
  }: { open: boolean; adjustmentCount: number } = $props();

  let confirming = $state(false);

  async function handleConfirm() {
    confirming = true;
    try {
      await confirmDelivery();
      toast.success("Entrega confirmada");
      open = false;
    } catch (err) {
      toast.error(String(err));
    } finally {
      confirming = false;
    }
  }
</script>

<AlertDialog.Root bind:open>
  <AlertDialog.Content>
    <AlertDialog.Header>
      <AlertDialog.Title>Confirmar entrega</AlertDialog.Title>
      <AlertDialog.Description>
        Se actualizará el saldo de <strong>{adjustmentCount}</strong>
        {adjustmentCount === 1 ? "registro" : "registros"}.
        ¿Deseas continuar?
      </AlertDialog.Description>
    </AlertDialog.Header>
    <AlertDialog.Footer>
      <AlertDialog.Cancel>Cancelar</AlertDialog.Cancel>
      <AlertDialog.Action onclick={handleConfirm} disabled={confirming}>
        {confirming ? "Confirmando..." : "Confirmar"}
      </AlertDialog.Action>
    </AlertDialog.Footer>
  </AlertDialog.Content>
</AlertDialog.Root>
