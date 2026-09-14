export namespace dto {
	
	export class DeliveryEntryDTO {
	    client_ref_uuid: string;
	    client_name: string;
	    client_address: string;
	    client_balance: number;
	    product_ref_uuid: string;
	    product_name: string;
	    charged_units: number;
	    credit_units: number;
	
	    static createFrom(source: any = {}) {
	        return new DeliveryEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_ref_uuid = source["client_ref_uuid"];
	        this.client_name = source["client_name"];
	        this.client_address = source["client_address"];
	        this.client_balance = source["client_balance"];
	        this.product_ref_uuid = source["product_ref_uuid"];
	        this.product_name = source["product_name"];
	        this.charged_units = source["charged_units"];
	        this.credit_units = source["credit_units"];
	    }
	}
	export class DeliveryHistoryEntryDTO {
	    client_ref_uuid: string;
	    client_name: string;
	    product_ref_uuid: string;
	    product_name: string;
	    charged_units: number;
	    credit_units: number;
	    delivery_date: string;
	
	    static createFrom(source: any = {}) {
	        return new DeliveryHistoryEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_ref_uuid = source["client_ref_uuid"];
	        this.client_name = source["client_name"];
	        this.product_ref_uuid = source["product_ref_uuid"];
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
	    client_ref_uuid: string;
	    client_name: string;
	    product_ref_uuid: string;
	    product_name: string;
	    charged_units: number;
	    credit_units: number;
	    delivery_date: string;
	
	    static createFrom(source: any = {}) {
	        return new ZoneReportEntryDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_ref_uuid = source["client_ref_uuid"];
	        this.client_name = source["client_name"];
	        this.product_ref_uuid = source["product_ref_uuid"];
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
	    zone_id: string;
	    create_date: string;
	    unregister_date?: string;
	    priority_order: number;
	    reference_uuid: string;
	
	    static createFrom(source: any = {}) {
	        return new Client(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.zone_id = source["zone_id"];
	        this.create_date = source["create_date"];
	        this.unregister_date = source["unregister_date"];
	        this.priority_order = source["priority_order"];
	        this.reference_uuid = source["reference_uuid"];
	    }
	}
	export class ClientProduct {
	    id: number;
	    client_id: string;
	    product_id: string;
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
	    reference_uuid: string;
	
	    static createFrom(source: any = {}) {
	        return new Product(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.type = source["type"];
	        this.active = source["active"];
	        this.reference_uuid = source["reference_uuid"];
	    }
	}
	export class Zone {
	    id: number;
	    name: string;
	    reference_uuid: string;
	
	    static createFrom(source: any = {}) {
	        return new Zone(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.reference_uuid = source["reference_uuid"];
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
	    client_ref_uuid: string;
	    product_ref_uuid: string;
	    delta: number;
	    charged_units: number;
	    delivery_date: string;
	
	    static createFrom(source: any = {}) {
	        return new BalanceAdjustment(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.client_ref_uuid = source["client_ref_uuid"];
	        this.product_ref_uuid = source["product_ref_uuid"];
	        this.delta = source["delta"];
	        this.charged_units = source["charged_units"];
	        this.delivery_date = source["delivery_date"];
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

