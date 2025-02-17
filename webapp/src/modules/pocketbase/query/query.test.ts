import { describe, it, expect, vi, beforeEach } from 'vitest';
import { PocketbaseQuery } from './query';
import type { CollectionName } from '@/pocketbase/collections-models';
import type PocketBase from 'pocketbase';

// Mock PocketBase client

const mockGetList = vi.fn();
const mockGetFullList = vi.fn();
const mockGetOne = vi.fn();

const mockPb = {
	collection: () => ({
		getList: mockGetList,
		getFullList: mockGetFullList,
		getOne: mockGetOne
	})
} as unknown as PocketBase;

//

describe('PocketbaseQuery', () => {
	beforeEach(() => {
		vi.clearAllMocks();
	});

	// Test basic query construction
	it('should create a basic query with default options', () => {
		const query = new PocketbaseQuery('users' as CollectionName, { pb: mockPb });
		expect(query.collection).toBe('users');
	});

	// Test pagination
	it('should handle pagination correctly', async () => {
		const query = new PocketbaseQuery('users' as CollectionName, { pb: mockPb }).withPagination(
			20
		);
		await query.getList(1);
		expect(mockGetList).toHaveBeenCalledWith(1, 20, expect.any(Object));
	});

	// Test filters
	it('should handle filters correctly', () => {
		const query = new PocketbaseQuery('users' as CollectionName, { pb: mockPb })
			.addFilter('age >= 18')
			.addFilter('verified = true');

		const options = query.buildOptions();
		expect(options.filter).toBe('age >= 18 && verified = true');
	});

	// Test search
	it('should handle search correctly', () => {
		const query = new PocketbaseQuery('users', { pb: mockPb }).addSearch('john', [
			'name',
			'email'
		]);

		const options = query.buildOptions();
		expect(options.filter).toBe("name ~ 'john' || email ~ 'john'");
	});

	// Test sort
	it('should handle sorting correctly', () => {
		const query = new PocketbaseQuery('users', { pb: mockPb })
			.addSort('created', 'DESC')
			.addSort('name', 'ASC');

		const options = query.buildOptions();
		expect(options.sort).toBe('-created,+name');
	});

	// Test flip sort
	it('should flip sort order correctly', () => {
		const query = new PocketbaseQuery('users', { pb: mockPb }).addSort('created', 'DESC');

		query.flipSort();
		const options = query.buildOptions();
		expect(options.sort).toBe('+created');
	});

	// Test exclude
	it('should handle exclude filters correctly', () => {
		const query = new PocketbaseQuery('users', { pb: mockPb }).addExclude(['123', '456']);

		const options = query.buildOptions();
		expect(options.filter).toBe("id != '123' && id != '456'");
	});

	// Test combined filters
	it('should handle combined filters correctly', () => {
		const query = new PocketbaseQuery('users', { pb: mockPb })
			.addFilter('age >= 18')
			.addSearch('john', ['name'])
			.addExclude(['123'])
			.addSort('created', 'DESC');

		const options = query.buildOptions();
		expect(options.filter).toBe("age >= 18 && name ~ 'john' && id != '123'");
		expect(options.sort).toBe('-created');
	});

	// Test API methods
	describe('API Methods', () => {
		it('should call getList with correct parameters', async () => {
			const query = new PocketbaseQuery('users', { pb: mockPb });
			await query.getList(1);
			expect(mockGetList).toHaveBeenCalledWith(1, 10, expect.any(Object));
		});

		it('should call getFullList with correct parameters', async () => {
			const query = new PocketbaseQuery('users', { pb: mockPb });
			await query.getFullList();
			expect(mockGetFullList).toHaveBeenCalledWith(expect.any(Object));
		});

		it('should call getOne with correct parameters', async () => {
			const query = new PocketbaseQuery('users', { pb: mockPb });
			await query.getOne('123');
			expect(mockGetOne).toHaveBeenCalledWith('123', expect.any(Object));
		});
	});

	// Test expand functionality
	it('should handle expand correctly', () => {
		const query = new PocketbaseQuery('users', {
			pb: mockPb,
			expand: ['authorizations_via_owner', 'z_test_collection_via_owner']
		});

		const options = query.buildOptions();
		expect(options.expand).toBe('authorizations_via_owner,z_test_collection_via_owner');
	});
});
