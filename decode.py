from __future__ import print_function
import sys
import os
import json
from WXBizDataCrypt import WXBizDataCrypt

pyPath = os.path.dirname(os.path.realpath(__file__))

def main(sessionKey, encryptedData, iv):
	config = open(os.path.join(pyPath, 'config.cfg'), 'r').read()
	appId = json.loads(config)['appid']
	# sessionKey = 'tiihtNczf5v6AKRyjwEUhQ=='
	# encryptedData = 'CiyLU1Aw2KjvrjMdj8YKliAjtP4gsMZMQmRzooG2xrDcvSnxIMXFufNstNGTyaGS9uT5geRa0W4oTOb1WT7fJlAC+oNPdbB+3hVbJSRgv+4lGOETKUQz6OYStslQ142dNCuabNPGBzlooOmB231qMM85d2/fV6ChevvXvQP8Hkue1poOFtnEtpyxVLW1zAo6/1Xx1COxFvrc2d7UL/lmHInNlxuacJXwu0fjpXfz/YqYzBIBzD6WUfTIF9GRHpOn/Hz7saL8xz+W//FRAUid1OksQaQx4CMs8LOddcQhULW4ucetDf96JcR3g0gfRK4PC7E/r7Z6xNrXd2UIeorGj5Ef7b1pJAYB6Y5anaHqZ9J6nKEBvB4DnNLIVWSgARns/8wR2SiRS7MNACwTyrGvt9ts8p12PKFdlqYTopNHR1Vf7XjfhQlVsAJdNiKdYmYVoKlaRv85IfVunYzO0IKXsyl7JCUjCpoG20f0a04COwfneQAGGwd5oa+T8yO5hzuyDb/XcxxmK01EpqOyuxINew=='
	# iv = 'r7BXXKkLb8qrSNn05n0qiA=='

	pc = WXBizDataCrypt(appId, sessionKey)

	print(json.dumps(pc.decrypt(encryptedData, iv)))

if __name__ == '__main__':
	f = open(sys.argv[1], 'r').readlines()
	main(f[0].strip(), f[1].strip(), f[2].strip())
