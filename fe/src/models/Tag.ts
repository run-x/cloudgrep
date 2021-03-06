export interface MockTag {
	Key: string;
	Value: string;
	Count?: number;
	ResourceIds?: string[];
}

export interface Tag {
	key: string;
	value: string;
	group: string;
}

export type TagType = { [key: string]: string };
