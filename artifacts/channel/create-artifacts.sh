# System channel
SYS_CHANNEL="sys-channel"

CHANNEL_NAME="bankachannel"
CHANNEL_PROFILE="BasicChannel"

# Delete existing artifacts
rm -rf ${CHANNEL_NAME}
rm -rf ../../channel-artifacts/*
mkdir ${CHANNEL_NAME}
#Generate Crypto artifactes for organizations
# cryptogen generate --config=./crypto-config.yaml --output=./crypto-config/

echo $CHANNEL_NAME

# Generate System Genesis block
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL -outputBlock ./${CHANNEL_NAME}/genesis-${CHANNEL_NAME}.block

cp ./${CHANNEL_NAME}/genesis-${CHANNEL_NAME}.block ./genesis.block

# Generate channel configuration block
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputCreateChannelTx ./${CHANNEL_NAME}/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME

echo "#######    Generating anchor peer update for Org1MSP  ##########"
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputAnchorPeersUpdate ./${CHANNEL_NAME}/Org1MSPanchors-${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

echo "#######    Generating anchor peer update for Org2MSP  ##########"
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputAnchorPeersUpdate ./${CHANNEL_NAME}/Org2MSPanchors-${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg Org2MSP

echo "#######    Generating anchor peer update for Org3MSP  ##########"
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputAnchorPeersUpdate ./${CHANNEL_NAME}/Org3MSPanchors-${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg Org3MSP

# ########################################################################################################################

CHANNEL_NAME="bankbchannel"
CHANNEL_PROFILE="SecondChannel"
mkdir ${CHANNEL_NAME}
echo $CHANNEL_NAME
# # Generate System Genesis block
configtxgen -profile OrdererGenesis -configPath . -channelID $SYS_CHANNEL -outputBlock ./${CHANNEL_NAME}/genesis-${CHANNEL_NAME}.block

# cp ./${CHANNEL_NAME}/genesis-${CHANNEL_NAME}.block ./${CHANNEL_NAME}/genesis.block

# Generate channel configuration block
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputCreateChannelTx ./${CHANNEL_NAME}/${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME

echo "#######    Generating anchor peer update for Org1MSP  ##########"
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputAnchorPeersUpdate ./${CHANNEL_NAME}/Org1MSPanchors-${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg Org1MSP

echo "#######    Generating anchor peer update for Org2MSP  ##########"
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputAnchorPeersUpdate ./${CHANNEL_NAME}/Org2MSPanchors-${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg Org2MSP

echo "#######    Generating anchor peer update for Org3MSP  ##########"
configtxgen -profile ${CHANNEL_PROFILE} -configPath . -outputAnchorPeersUpdate ./${CHANNEL_NAME}/Org3MSPanchors-${CHANNEL_NAME}.tx -channelID $CHANNEL_NAME -asOrg Org3MSP
