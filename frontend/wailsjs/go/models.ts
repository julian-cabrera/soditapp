export namespace dto {
	
	export class DeliveryEntryDTO {
	    client_id: number;
	    client_name: string;
	    client_address: string;
	    client_balance: number;
	    product_id: number;
	    product_name: string;
	    charged_units: number;
	    credit_units: number;
	
	    static createFrom(source: any = {}) {
	        return new DeliveryEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_id = source["client_id"];
	        this.client_name = source["client_name"];
	        this.client_address = source["client_address"];
	        this.client_balance = source["client_balance"];
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.charged_units = source["charged_units"];
	        this.credit_units = source["credit_units"];
	    }
	}
	export class DeliveryHistoryEntryDTO {
	    client_id: number;
	    client_name: string;
	    product_id: number;
	    product_name: string;
	    charged_units: number;
	    credit_units: number;
	    delivery_date: string;
	
	    static createFrom(source: any = {}) {
	        return new DeliveryHistoryEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_id = source["client_id"];
	        this.client_name = source["client_name"];
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.charged_units = source["charged_units"];
	        this.credit_units = source["credit_units"];
	        this.delivery_date = source["delivery_date"];
	    }
	}
	export class DeliverySheetDTO {
	    zone_id: number;
	    zone_name: string;
	    // Go type: time
	    created_at: any;
	    page_number: number;
	    entries: DeliveryEntryDTO[];
	
	    static createFrom(source: any = {}) {
	        return new DeliverySheetDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.zone_id = source["zone_id"];
	        this.zone_name = source["zone_name"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.page_number = source["page_number"];
	        this.entries = this.convertValues(source["entries"], DeliveryEntryDTO);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ZoneReportEntryDTO {
	    client_id: number;
	    client_name: string;
	    product_id: number;
	    product_name: string;
	    charged_units: number;
	    credit_units: number;
	    delivery_date: string;
	
	    static createFrom(source: any = {}) {
	        return new ZoneReportEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_id = source["client_id"];
	        this.client_name = source["client_name"];
	        this.product_id = source["product_id"];
	        this.product_name = source["product_name"];
	        this.charged_units = source["charged_units"];
	        this.credit_units = source["credit_units"];
	        this.delivery_date = source["delivery_date"];
	    }
	}

}

export namespace model {
	
	export class Client {
	    id: number;
	    name: string;
	    address: string;
	    zone_id: number;
	    // Go type: time
	    create_date: any;
	    // Go type: time
	    unregister_date?: any;
	    priority_order: number;
	
	    static createFrom(source: any = {}) {
	        return new Client(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.zone_id = source["zone_id"];
	        this.create_date = this.convertValues(source["create_date"], null);
	        this.unregister_date = this.convertValues(source["unregister_date"], null);
	        this.priority_order = source["priority_order"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class ClientProduct {
	    id: number;
	    client_id: number;
	    product_id: number;
	    balance: number;
	
	    static createFrom(source: any = {}) {
	        return new ClientProduct(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.client_id = source["client_id"];
	        this.product_id = source["product_id"];
	        this.balance = source["balance"];
	    }
	}
	export class Product {
	    id: number;
	    name: string;
	    type: string;
	    active: boolean;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.active = source["active"];
	    }
	}
	export class Zone {
	    id: number;
	    name: string;
	
	    static createFrom(source: any = {}) {
	        return new Zone(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	    }
	}

}

export namespace repository {
	
	export class ClientProductRepository {
	
	
	    static createFrom(source: any = {}) {
	        return new ClientProductRepository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ClientRepository {
	
	
	    static createFrom(source: any = {}) {
	        return new ClientRepository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class DeliveryTransactionRepository {
	
	
	    static createFrom(source: any = {}) {
	        return new DeliveryTransactionRepository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ProductRepository {
	
	
	    static createFrom(source: any = {}) {
	        return new ProductRepository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class ZoneRepository {
	
	
	    static createFrom(source: any = {}) {
	        return new ZoneRepository(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

export namespace service {
	
	export class BalanceAdjustment {
	    client_id: number;
	    product_id: number;
	    delta: number;
	    charged_units: number;
	    delivery_date: string;
	
	    static createFrom(source: any = {}) {
	        return new BalanceAdjustment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_id = source["client_id"];
	        this.product_id = source["product_id"];
	        this.delta = source["delta"];
	        this.charged_units = source["charged_units"];
	        this.delivery_date = source["delivery_date"];
	    }
	}
	export class ImportResult {
	    inserted: number;
	    skipped: number;
	    errors: string[];
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.inserted = source["inserted"];
	        this.skipped = source["skipped"];
	        this.errors = source["errors"];
	    }
	}

}

export namespace sql {
	
	export class DB {
	
	
	    static createFrom(source: any = {}) {
	        return new DB(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

