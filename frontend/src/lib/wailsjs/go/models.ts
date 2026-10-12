export namespace main {
	
	export class DocumentDTO {
	    id: string;
	    title: string;
	    filePath: string;
	    content: string;
	
	    static createFrom(source: any = {}) {
	        return new DocumentDTO(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.filePath = source["filePath"];
	        this.content = source["content"];
	    }
	}

}

