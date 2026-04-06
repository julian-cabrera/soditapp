<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Input } from "$lib/components/ui/input/index.js";
  import { Label } from "$lib/components/ui/label/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { createProduct, updateProduct } from "$lib/stores/products.svelte.js";
  import { toast } from "svelte-sonner";
  import type { model } from "wailsjs/go/models";

  let { open = $bindable(false), product }: { open: boolean; product?: model.Product } = $props();

  let name = $state("");
  let type = $state("");
  let submitting = $state(false);

  let isEdit = $derived(!!product);
  let title = $derived(isEdit ? "Renombrar Producto" : "Nuevo Producto");

  $effect(() => {
    if (open) {
      name = product?.name ?? "";
      type = product?.type ?? "";
    }
  });

  async function handleSubmit(e: SubmitEvent) {
    e.preventDefault();
    const trimmedName = name.trim();
    const trimmedType = type.trim();
    if (!trimmedName || !trimmedType) return;

    submitting = true;
    try {
      if (isEdit && product) {
        await updateProduct(product.id, trimmedName);
        toast.success("Producto actualizado");
      } else {
        await createProduct(trimmedName, trimmedType);
        toast.success("Producto creado");
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
        <Label for="product-name">Nombre</Label>
        <Input id="product-name" bind:value={name} placeholder="Nombre del producto" required />
      </div>
      <div class="space-y-2">
        <Label for="product-type">Tipo</Label>
        <Input id="product-type" bind:value={type} placeholder="Tipo de producto" required />
      </div>
      <Dialog.Footer>
        <Button variant="outline" type="button" onclick={() => (open = false)}>Cancelar</Button>
        <Button type="submit" disabled={submitting || !name.trim() || !type.trim()}>
          {submitting ? "Guardando..." : "Guardar"}
        </Button>
      </Dialog.Footer>
    </form>
  </Dialog.Content>
</Dialog.Root>
