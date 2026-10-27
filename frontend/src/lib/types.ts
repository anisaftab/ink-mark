/** Mirrors DocumentDTO from app.go */
export interface Document {
	id: string;
	title: string;
	filePath: string;
	content: string;
	favorited: boolean;
}

/** Mirrors ProjectDTO from app.go */
export interface Project {
	id: string;
	name: string;
	folderPath: string;
	documents: Document[];
}
