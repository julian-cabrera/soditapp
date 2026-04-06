import {
  CreateProduct,
  ListProducts,
  UpdateProduct,
  ToggleProductActive,
} from "wailsjs/go/service/ProductService";
import type { model } from "wailsjs/go/models";

let products = $state<model.Product[]>([]);
let loading = $state(false);

export function getProducts() {
  return products;
}

export function getProductsLoading() {
  return loading;
}

export async function fetchProducts() {
  loading = true;
  try {
    products = await ListProducts();
  } finally {
    loading = false;
  }
}

export async function createProduct(name: string, type: string) {
  loading = true;
  try {
    await CreateProduct(name, type);
    await fetchProducts();
  } finally {
    loading = false;
  }
}

export async function updateProduct(id: number, name: string) {
  loading = true;
  try {
    await UpdateProduct(id, name);
    await fetchProducts();
  } finally {
    loading = false;
  }
}

export async function toggleProductActive(id: number) {
  loading = true;
  try {
    await ToggleProductActive(id);
    await fetchProducts();
  } finally {
    loading = false;
  }
}
