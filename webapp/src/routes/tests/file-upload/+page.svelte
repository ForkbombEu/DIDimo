<script lang="ts">
	import { pb } from '@/pocketbase';
	import { createDummyFile } from '@/utils/other';

	pb.collection('z_test_collection')
		.update('bn0sadak3y1876r', {
			// file_multi_field: createDummyFile()
			'file_multi_field-': ['file_3I6Ivkk46N.txt']
		})
		.then((e) => console.log(e))
		.catch((e) => console.log(e));
</script>
