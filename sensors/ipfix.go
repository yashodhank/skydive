/*
 * Copyright (C) 2015 Red Hat, Inc.
 *
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 *
 */

package sensors

import (
	"bytes"
	"fmt"
	"net"

	"github.com/calmh/ipfix"
)

type IpfixSensor struct {
	Addr string
	Port int
}

func (sensor *IpfixSensor) Start() error {
	var buf [4096]byte

	addr := net.UDPAddr{
		Port: sensor.Port,
		IP:   net.ParseIP(sensor.Addr),
	}
	conn, err := net.ListenUDP("udp", &addr)
	defer conn.Close()
	if err != nil {
		return err
	}

	s := ipfix.NewSession()
	i := ipfix.NewInterpreter(s)
	_ = i

	reader := bytes.NewReader(buf[:])
	for {
		_, _, err := conn.ReadFromUDP(buf[:])
		// ParseReader will block until a full message is available.
		msg, err := s.ParseReader(reader)
		if err != nil {
			panic(err)
		}
		reader.Seek(0, 0)
		fmt.Println("-----------------------------")
		fmt.Println(msg)
		fmt.Println("-----------------------------")

		for _, record := range msg.DataRecords {
			fmt.Println(i.Interpret(record))
			// record contains raw enterpriseId, fieldId => []byte information
			//fmt.Println(record)
		}
	}

	return nil
}

func NewIpfixSensor(addr string, port int) IpfixSensor {
	sensor := IpfixSensor{Addr: addr, Port: port}
	return sensor
}