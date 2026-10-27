export namespace main {
	
	export class DocumentDTO {
	    id: string;
	    title: string;
	    filePath: string;
	    content: string;
	    favorited: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DocumentDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.filePath = source["filePath"];
	        this.content = source["content"];
	        this.favorited = source["favorited"];
	    }
	}
	export class ProjectDTO {
	    id: string;
	    name: string;
	    folderPath: string;
	    documents: DocumentDTO[];
	
	    static createFrom(source: any = {}) {
	        return new ProjectDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.folderPath = source["folderPath"];
	        this.documents = this.convertValues(source["documents"], DocumentDTO);
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

}

