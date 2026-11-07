<script lang="ts">
  import * as Dialog from "$lib/components/ui/dialog/index.js";
  import { Button } from "$lib/components/ui/button/index.js";
  import { HugeiconsIcon } from "@hugeicons/svelte";
  import { PlusSignIcon, RotateLeftIcon } from "@hugeicons/core-free-icons";

  let {
    open = $bindable(false),
    onAdd,
    onRestore,
  }: { open: boolean; onAdd: () => void; onRestore: () => void } = $props();

  function handleAdd() {
    onAdd();
    open = false;
  }

  function handleRestore() {
    onRestore();
    open = false;
  }
</script>

<Dialog.Root bind:open>
  <Dialog.Content class="sm:max-w-sm">
    <Dialog.Header>
      <Dialog.Title>Nuevo Cliente</Dialog.Title>
      <Dialog.Description>
        ¿Deseas agregar un nuevo cliente o restaurar uno previamente dado de
        baja?
      </Dialog.Description>
    </Dialog.Header>
    <div class="flex flex-col gap-3">
      <Button size="lg" class="w-full" onclick={handleAdd}>
        <HugeiconsIcon icon={PlusSignIcon} size={16} />
        Agregar Nuevo
      </Button>
      <Button
        size="lg"
        variant="secondary"
        class="w-full"
        onclick={handleRestore}
      >
        <HugeiconsIcon icon={RotateLeftIcon} size={16} />
        Restaurar Anterior
      </Button>
      <Button variant="outline" class="w-full" onclick={() => (open = false)}
        >Cancelar</Button
      >
    </div>
  </Dialog.Content>
</Dialog.Root>
