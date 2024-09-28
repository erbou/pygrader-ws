#!/bin/bash

main() {
        # Prepare the folders and db files for openssl ca
	set -eCu
	declare dir=./ca
	for p in root intr; do
		mkdir -m 0700 -p ${dir}/${p}/{certs,crl,newcerts,private,csr}
		echo 01   >& /dev/null > ${dir}/${p}/serial || true
		echo 0100 >& /dev/null > ${dir}/${p}/crlnumber || true
		touch ${dir}/${p}/index.txt
	done
	mkdir -p ${dir}/serv/{certs,crl,csr,private}
	chmod -R u+rwX,og-rwx ./ca
}

main "${@}"
