package substrate

func (sh *SubstrateHarvester) GetActiveEra() (uint32, error) {

	activeEraKey, err := sh.GetStorageDataKey("Staking", "ActiveEra")
	if err != nil {
		return 0, err
	}

	var activeEra uint32
	err = sh.GetStorageLatest(activeEraKey, &activeEra)
	if err != nil {
		return 0, err
	}
	return activeEra, nil

}

func (sh *SubstrateHarvester) GetActiveEraDepth() ([]byte, error) {
	activeEra, err := sh.GetActiveEra()
	if err != nil {
		return nil, err
	}
	return i32tob(activeEra), nil
}

func (sh *SubstrateHarvester) GetEraDepth(era uint32) ([]byte, error) {
	return i32tob(era), nil
}

func i32tob(val uint32) []byte {
	r := make([]byte, 4)
	for i := uint32(0); i < 4; i++ {
		r[i] = byte((val >> (8 * i)) & 0xff)
	}
	return r
}
