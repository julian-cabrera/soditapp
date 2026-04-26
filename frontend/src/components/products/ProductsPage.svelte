<script lang="ts">
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { PencilEdit01Icon, PlusSignIcon, ToggleOnIcon, ToggleOffIcon } from "@hugeicons/core-free-icons";
  import { Button } from "$lib/components/ui/button/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import { Badge } from "$lib/components/ui/badge/index.js";
  import {
    getProducts,
    getProductsLoading,
    fetchProducts,
    toggleProductActive,
  } from "$lib/stores/products.svelte.js";
  import ProductFormDialog from "./ProductFormDialog.svelte";
  import { toast } from "svelte-sonner";
  import type { model } from "wailsjs/go/models";

  let formOpen = $state(false);
  let selectedProduct = $state<model.Product | undefined>(undefined);

  let products = $derived(getProducts());
  let loading = $derived(getProductsLoading());

  $effect(() => {
    fetchProducts();
  });

  function openCreate() {
    selectedProduct = undefined;
    formOpen = true;
  }

  function openEdit(product: model.Product) {
    selectedProduct = product;
    formOpen = true;
  }

  async function handleToggle(product: model.Product) {
    try {
      await toggleProductActive(product.id);
      toast.success(product.active ? "Producto desactivado" : "Producto activado");
    } catch (err) {
      toast.error(String(err));
    }
  }
</script>

<div class="space-y-4">
  <div class="flex items-center justify-between">
    <h2 class="text-lg font-semibold">Productos</h2>
    <Button size="default" onclick={openCreate}>
      <HugeiconsIcon icon={PlusSignIcon} size={14} />
      Nuevo Producto
    </Button>
  </div>

  {#if loading && products.length === 0}
    <p class="text-sm text-muted-foreground">Cargando...</p>
  {:else if products.length === 0}
    <p class="text-sm text-muted-foreground">No hay productos registrados.</p>
  {:else}
    <Table.Root>
      <Table.Header>
        <Table.Row>
          <Table.Head>Nombre</Table.Head>
          <Table.Head>Tipo</Table.Head>
          <Table.Head>Estado</Table.Head>
          <Table.Head class="w-32 text-right">Acciones</Table.Head>
        </Table.Row>
      </Table.Header>
      <Table.Body>
        {#each products as product (product.id)}
          <Table.Row>
            <Table.Cell>{product.name}</Table.Cell>
            <Table.Cell>{product.type}</Table.Cell>
            <Table.Cell>
              {#if product.active}
                <Badge variant="default">Activo</Badge>
              {:else}
                <Badge variant="secondary">Inactivo</Badge>
              {/if}
            </Table.Cell>
            <Table.Cell class="text-right">
              <div class="flex justify-end gap-1">
                <Button variant="ghost" size="icon" onclick={() => openEdit(product)} title="Renombrar">
                  <HugeiconsIcon icon={PencilEdit01Icon} size={14} />
                </Button>
                <Button
                  variant="ghost"
                  size="icon"
                  onclick={() => handleToggle(product)}
                  title={product.active ? "Desactivar" : "Activar"}
                >
                  <HugeiconsIcon icon={product.active ? ToggleOnIcon : ToggleOffIcon} size={14} />
                </Button>
              </div>
            </Table.Cell>
          </Table.Row>
        {/each}
      </Table.Body>
    </Table.Root>
  {/if}
</div>

<ProductFormDialog bind:open={formOpen} product={selectedProduct} />
