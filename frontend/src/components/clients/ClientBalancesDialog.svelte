<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import * as Table from "$lib/components/ui/table/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { getProducts, fetchProducts } from "$lib/stores/products.svelte.js";
  import { GetClientBalances } from "wailsjs/go/service/ClientProductService";
  import { toast } from "svelte-sonner";
  import type { model } from "wailsjs/go/models";

  let { open = $bindable(false), client }: { open: boolean; client: model.Client } = $props();

  let balances = $state<model.ClientProduct[]>([]);
  let loading = $state(false);

  let products = $derived(getProducts());
  let productMap = $derived(new Map(products.map((p) => [p.id, p.name])));

  $effect(() => {
    if (open) {
      loadBalances();
    }
  });

  async function loadBalances() {
    loading = true;
    try {
      await fetchProducts();
      balances = await GetClientBalances(client.id);
    } catch (err) {
      toast.error(String(err));
    } finally {
      loading = false;
    }
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="sm:max-w-md">
    <Dialog.Header>
      <Dialog.Title>Saldos de {client.name}</Dialog.Title>
    </Dialog.Header>

    {#if loading}
      <p class="text-sm text-muted-foreground">Cargando...</p>
    {:else if balances.length === 0}
      <p class="text-sm text-muted-foreground">No hay saldos registrados.</p>
    {:else}
      <Table.Root>
        <Table.Header>
          <Table.Row>
            <Table.Head>Producto</Table.Head>
            <Table.Head class="text-right">Saldo</Table.Head>
          </Table.Row>
        </Table.Header>
        <Table.Body>
          {#each balances as cp (cp.id)}
            <Table.Row>
              <Table.Cell>{productMap.get(cp.product_id) ?? "—"}</Table.Cell>
              <Table.Cell class="text-right">{cp.balance}</Table.Cell>
            </Table.Row>
          {/each}
        </Table.Body>
      </Table.Root>
    {/if}

    <Dialog.Footer>
      <Button variant="outline" onclick={() => (open = false)}>Cerrar</Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
